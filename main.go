package main

import (
	_ "embed"
	"log"
	"net/http"
	"netdash/internal/handler"
	"netdash/internal/logger"
	"netdash/internal/repository"
	"netdash/internal/server"
	"netdash/internal/service"
	"netdash/internal/utils"
	"netdash/web"
	"strings"
)

//go:embed VERSION
var versionFile string

func main() {
	version := strings.TrimSpace(versionFile)
	version = strings.TrimPrefix(version, "v")

	if version == "" {
		version = "unknown"
	}

	log.SetFlags(0)
	logger.Log("SYSTEM", "Starting NetDash v%s", version)

	repo := repository.NewSQLiteRepository()
	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}

	speedService := service.NewSpeedtestService(repo)
	schedulerService := service.NewSchedulerService(repo, speedService)
	schedulerService.Start()

	appHandler := handler.NewHandler(repo, speedService, schedulerService, web.Assets)

	e := server.New(appHandler, web.Assets, version)

	utils.LogAccessInfo(":80")
	if err := e.Start(":80"); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}