package server

import (
	"time"

	"github.com/labstack/echo/v4"
	configpb "github.com/lynx-go/lynx-template/internal/pkg/config"
	"github.com/lynx-go/lynx/server/http"
)

func NewHTTPServer(
	config *configpb.AppConfig,
	router *echo.Echo,
) *http.Server {
	c := config.GetServer().GetHttp()
	addr := c.Addr
	timeout := parseTimeout(c.Timeout)

	return http.NewServer(router, http.WithAddr(addr), http.WithTimeout(timeout))
}

func parseTimeout(s string) time.Duration {
	if s == "" {
		s = "60s"
	}
	timeout, _ := time.ParseDuration(s)
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	return timeout
}
