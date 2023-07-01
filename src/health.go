package src

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthStatus struct {
	Status string `json:"status"`
}

func init() {
	GetEchoRoot().GET("health", health)
}

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, &HealthStatus{Status: "ok"})
}
