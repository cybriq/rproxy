package main

import (
	"github.com/labstack/echo/v4"
	"net"
)

type (
	Host struct {
		*echo.Echo
	}
)

const (
	listenHost = "localhost"
	listenPort = "8080"
)

func main() {
	// Hosts
	hosts := map[string]*Host{
		"api." + listenHost:    {getEcho("API")},
		"static." + listenHost: {getEcho("Static assets")},
		"blog." + listenHost:   {getEcho("Blog")},
		listenHost:             {getEcho("Landing page")},
	}
	// Server
	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		h, _, _ := net.SplitHostPort(req.Host)
		host := hosts[h]
		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}
		return
	},
	)
	e.Logger.Fatal(e.Start(":" + listenPort))
}
