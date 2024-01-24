package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Heartbeat(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}