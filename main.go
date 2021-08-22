package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/url"
	"strings"
)

type (
	proxySpec struct {
		path   string
		target *url.URL
	}
)

const (
	listenHost = "cybriq.systems"
	listenPort = "443"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

var (
	proxies = []proxySpec{
		{"/", mustParseURL("http://localhost:8001")},
		{"/git", mustParseURL("http://localhost:8002")},
		{"/blog", mustParseURL("http://localhost:8003")},
		{"/api", mustParseURL("http://localhost:8004")},
	}
)

func main() {
	// set up automatic redirect to https
	es := echo.New()
	es.Pre(middleware.HTTPSRedirect())
	go es.Start(listenHost + ":80")
	defer es.Close()
	// set up https enabled reverse proxy
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
						Name: proxies[i].path,
						URL:  proxies[i].target,
						Meta: nil,
					},
				},
			),
		),
		)
		hosts[proxies[i].path] = &Host{ec}
	}
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		h := req.URL.Path
		// fmt.Println(h)
		// if the hostname is the non-path this reverse proxy will apply
		host := hosts["/"]
		for i := range proxies {
			pth := proxies[i].path
			if pth == "/" {
				continue
			}
			if strings.HasPrefix(h, pth) {
				// fmt.Println(pth)
				host = hosts[pth]
				break
			}
		}
		host.Echo.ServeHTTP(res, req)
		return
	},
	)
	e.Logger.Fatal(e.StartTLS(listenHost+":"+listenPort, "/home/loki/cybriq_systems/deploy/cybriq_systems.crt",
		"/home/loki/cybriq_systems/deploy/cybriq_systems.key",
	),
	)
}

func mustParseURL(u string) *url.URL {
	url1, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return url1
}
