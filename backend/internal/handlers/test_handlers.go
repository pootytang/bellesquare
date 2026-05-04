package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Ping(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Pong!",
		"status":  "success",
		"db":      "connected",
	})
}