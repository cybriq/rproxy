package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func getEcho(echoText string) (ech *echo.Echo) {
	ech = echo.New()
	ech.Use(middleware.Logger())
	ech.Use(middleware.Recover())

	ech.GET("/",
		func(c echo.Context) error {
			return c.String(http.StatusOK, echoText)
		},
	)
	return
}
