package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = modifyResponse()
	proxy.Director = director()
	return proxy, nil
}

func director() func(*http.Request) {
	return func(req *http.Request) {
		// spew.Dump(req)
		req.URL.Scheme = "http"
		// req.URL.Host = "localhost:3000"
		// fmt.Println(req.URL.String())
		if strings.HasPrefix(req.URL.Path, "/git") {
			req.URL.Host = "localhost:3000"
			fmt.Println(req.URL.Path)
			split := strings.Split(req.URL.Path, "/git")
			if len(split) == 2 {
				req.URL.Path = strings.Join(split[1:], "")+"/"
			}
			fmt.Println(req.URL)
			// req.URL.Path
		}
		
		// if strings.HasPrefix(req.URL.Path, "/") {
		// 	req.URL.Host = target2.Host
		// }
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {

		// resp.Header.Set("X-Proxy", "Magical")
		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	
	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServeTLS(
		":443",
		"/home/loki/cybriq_systems/deploy/cybriq_systems.crt",
		"/home/loki/cybriq_systems/deploy/cybriq_systems.key",
		nil,
	),
	)
}
