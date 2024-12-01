package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *application) registerHealthCheckRoutes(e *echo.Group) {
	e.GET("/ping", app.ping)
}

func (app *application) ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
