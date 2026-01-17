package handler

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"netdash/internal/model"
	"netdash/internal/repository"
	"netdash/internal/service"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

type TemplateRegistry struct {
	templates *template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Handler struct {
	repo      repository.Repository
	speed     service.SpeedtestService
	scheduler service.SchedulerService
	assets    fs.FS
}

func NewHandler(r repository.Repository, s service.SpeedtestService, sch service.SchedulerService, assets fs.FS) *Handler {
	return &Handler{
		repo:      r,
		speed:     s,
		scheduler: sch,
		assets:    assets,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	tmpl, err := template.ParseFS(h.assets, "template/*.html", "template/partials/*.html", "template/layout/*.html")
	if err != nil {
		panic(err)
	}

	t := &TemplateRegistry{
		templates: tmpl,
	}
	e.Renderer = t

	e.GET("/", h.Dashboard)
	e.GET("/partials/status", h.GetStatus)
	e.GET("/partials/history", h.GetHistory)
	e.GET("/settings", h.SettingsPage)
	e.POST("/settings", h.SaveSettings)
	e.POST("/settings/clear", h.ClearDatabase)
	e.POST("/run-test", h.RunTestManual)
	e.GET("/api/stats", h.GetStats)
	e.DELETE("/delete/:id", h.DeleteResult)
}

func (h *Handler) Dashboard(c echo.Context) error {
	conf, _ := h.repo.GetConfig()
	limit := 10
	if conf != nil && conf.HistoryLimit > 0 {
		limit = conf.HistoryLimit
	}

	results, _ := h.repo.GetLatestResults(limit)
	stats, _ := h.repo.GetDBStats()
	version := c.Get("AppVersion")

	var latest interface{}
	if len(results) > 0 {
		latest = results[0]
	}
	
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Results":    results,
		"Latest":     latest,
		"IsRunning":  h.speed.IsRunning(),
		"Stats":      stats,
		"AppVersion": version,
		"Limit":      limit,
	})
}

func (h *Handler) GetStatus(c echo.Context) error {
	return c.Render(http.StatusOK, "status_button.html", map[string]interface{}{
		"IsRunning": h.speed.IsRunning(),
	})
}

func (h *Handler) GetHistory(c echo.Context) error {
	conf, _ := h.repo.GetConfig()
	limit := 10
	if conf != nil && conf.HistoryLimit > 0 {
		limit = conf.HistoryLimit
	}

	results, _ := h.repo.GetLatestResults(limit)
	for _, res := range results {
		if err := c.Render(http.StatusOK, "result_row.html", res); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) SettingsPage(c echo.Context) error {
	conf, _ := h.repo.GetConfig()
	stats, _ := h.repo.GetDBStats()
	version := c.Get("AppVersion")

	return c.Render(http.StatusOK, "settings.html", map[string]interface{}{
		"Config":     conf,
		"Stats":      stats,
		"AppVersion": version,
	})
}

func (h *Handler) SaveSettings(c echo.Context) error {
	serverID := c.FormValue("server_id")
	schedule := c.FormValue("cron_schedule")
	limitStr := c.FormValue("history_limit")

	if schedule != "manual" && schedule != "" {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err := parser.Parse(schedule); err != nil {
			return c.HTML(http.StatusBadRequest, "Invalid Cron Schedule")
		}
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	h.repo.UpdateConfig(&model.AppConfig{
		OoklaServerID: serverID,
		CronSchedule:  schedule,
		HistoryLimit:  limit,
	})
	
	h.scheduler.UpdateSchedule()
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) ClearDatabase(c echo.Context) error {
	h.repo.ClearResults()
	return c.Redirect(http.StatusSeeOther, "/settings")
}

func (h *Handler) RunTestManual(c echo.Context) error {
	res, err := h.speed.Run()
	if err != nil {
		if err.Error() == "test is already running" {
			return c.NoContent(http.StatusOK)
		}
		return c.HTML(http.StatusOK, "") 
	}

	c.Response().Header().Set("HX-Trigger", "refreshChart")
	return c.Render(http.StatusOK, "result_row.html", res)
}

func (h *Handler) GetStats(c echo.Context) error {
	conf, _ := h.repo.GetConfig()
	limit := 10
	if conf != nil && conf.HistoryLimit > 0 {
		limit = conf.HistoryLimit
	}

	results, err := h.repo.GetGraphData(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, results)
}

func (h *Handler) DeleteResult(c echo.Context) error {
	id := c.Param("id")
	h.repo.DeleteResult(id)
	return c.NoContent(http.StatusOK)
}