package server

import (
	"crypto/subtle"
	"io/fs"
	"log"
	"netdash/internal/handler"
	"netdash/internal/logger"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(h *handler.Handler, assets fs.FS, version string) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("AppVersion", version)
			return next(c)
		}
	})

	setupAuth(e)

	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		log.Fatal(err)
	}
	e.StaticFS("/static", staticFS)

	h.RegisterRoutes(e)

	return e
}

func setupAuth(e *echo.Echo) {
	user := os.Getenv("NETDASH_USER")
	pass := os.Getenv("NETDASH_PASSWORD")

	if user != "" && pass != "" {
		logger.Log("SECURITY", "Basic Auth ENABLED (User: %s)", user)
		e.Use(middleware.BasicAuth(func(u, p string, c echo.Context) (bool, error) {
			if subtle.ConstantTimeCompare([]byte(u), []byte(user)) == 1 &&
				subtle.ConstantTimeCompare([]byte(p), []byte(pass)) == 1 {
				return true, nil
			}
			return false, nil
		}))
	} else if user != "" || pass != "" {
		logger.Log("SECURITY", "WARNING: Partial Auth Configuration Detected!")
		logger.Log("SECURITY", "Missing NETDASH_USER or NETDASH_PASSWORD")
		logger.Log("SECURITY", "Authentication is DISABLED")
	} else {
		logger.Log("SECURITY", "No Auth Configured (Public Access Mode)")
	}
}