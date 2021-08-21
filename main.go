package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net"
	"net/url"
	"strings"
)

type (
	proxySpec struct {
		subdomain string
		target    *url.URL
	}
)

const (
	listenHost = "localhost"
	listenPort = "8080"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

var (
	proxies = []proxySpec{
		{"", mustParseURL("http://localhost:8001")},
		{"git", mustParseURL("http://localhost:8002")},
		{"blog", mustParseURL("http://localhost:8003")},
		{"api", mustParseURL("http://localhost:8004")},
	}
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	// set up reverse proxies
	hosts := map[string]*Host{}
	for i := range proxies {
		ec := echo.New()
		ec.Use(middleware.Logger())
		ec.Use(middleware.Recover())
		ec.Use(middleware.Proxy(
			middleware.NewRandomBalancer(
				[]*middleware.ProxyTarget{
					{
						Name: proxies[i].subdomain,
						URL:  proxies[i].target,
						Meta: nil,
					},
				},
			),
		),
		)
		hosts[proxies[i].subdomain] = &Host{ec}
	}
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		h, _, _ := net.SplitHostPort(req.Host)
		// if the hostname is the non-subdomain this reverse proxy will apply
		host := hosts[""]
		if h != listenHost {
			fmt.Println("host", h)
			for i := range proxies {
				prefix := proxies[i].subdomain
				if prefix == "" {
					continue
				}
				fmt.Println("prefix", prefix)
				if strings.HasPrefix(h, prefix) {
					host = hosts[proxies[i].subdomain]
					break
				}
			}
		}
		host.Echo.ServeHTTP(res, req)
		return
	},
	)
	e.Logger.Fatal(e.Start(":" + listenPort))
}

func mustParseURL(u string) *url.URL {
	url1, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url1
}
