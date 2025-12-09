package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lynx-go/lynx-app-template/internal/api/http"
	"github.com/lynx-go/lynx-app-template/pkg/jsonapi"
)

func NewRouter(helloApi *http.HelloAPI) *echo.Echo {
	root := newRouter()
	{
		root.POST("/api/hello", jsonapi.H(helloApi.Hello))
	}

	return root
}

func newRouter() *echo.Echo {
	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	return r
}
