package main

import (
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerUrls struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

func ProxyHandler(target string) echo.HandlerFunc {
	return func(c echo.Context) error {
		targetURL, err := url.Parse(target)
		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		c.Request().URL.Host = targetURL.Host
		c.Request().URL.Scheme = targetURL.Scheme
		c.Request().Host = targetURL.Host

		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}

func ShadowProxyHandler(target string) {
	go func(target string) {
		targetURL, err := url.Parse(target)
		if err != nil {
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Create a dummy request and response recorder
		req := httptest.NewRequest("GET", targetURL.String(), nil)
		resp := httptest.NewRecorder()

		req = req.Clone(req.Context())

		proxy.ServeHTTP(resp, req)
	}(target)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	servers := []ServerUrls{
		{
			URL:  "http://localhost:3000",
			Type: "prod",
		},
		{
			URL:  "http://localhost:4000",
			Type: "shadow",
		},
	}

	var production string
	var shadow string

	for _, server := range servers {
		if server.Type == "prod" {
			production = server.URL
		} else if server.Type == "shadow" {
			shadow = server.URL
		}
	}

	e.Any("/*", func(c echo.Context) error {
		ShadowProxyHandler(shadow)
		return ProxyHandler(production)(c)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
