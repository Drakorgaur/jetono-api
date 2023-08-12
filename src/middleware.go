package src

import (
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	e = GetEchoRoot()

	// Middleware
	// TODO
	e.Use(resetCliMiddleware())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
