package server

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lynx-go/lynx-app-template/internal/api/http"
	configpb "github.com/lynx-go/lynx-app-template/internal/pkg/config"
	"github.com/lynx-go/lynx-app-template/pkg/echoutil"
	"github.com/lynx-go/lynx-app-template/pkg/jsonapi"
)

func NewRouter(
	slogger *slog.Logger,
	config *configpb.AppConfig,
	helloApi *http.HelloAPI,
	authApi *http.AccountAPI,
) *echo.Echo {
	jwtSecret := config.GetSecurity().GetJwt().GetSecret()

	r := echoutil.NewEcho(slogger)

	middlewares(r, config)

	publicGroup := r.Group("")
	{
		publicGroup.POST("/sign-up", jsonapi.H(authApi.SignUp))
		publicGroup.POST("/token", authApi.Token)
	}
	secureGroup := r.Group("", echoutil.JWTMiddleware(jwtSecret))
	{
		secureGroup.POST("/logout", authApi.Logout)
	}
	apiGroup := secureGroup.Group("/api")
	{
		apiGroup.POST("/hello", jsonapi.H(helloApi.Hello))
	}

	return r
}

func middlewares(r *echo.Echo, config *configpb.AppConfig) {
	c := config.GetServer().GetHttp()
	if c.Cors != nil {
		r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     c.Cors.AllowOrigins,
			AllowHeaders:     c.Cors.AllowHeaders,
			AllowMethods:     c.Cors.AllowMethods,
			ExposeHeaders:    c.Cors.ExposeHeaders,
			AllowCredentials: c.Cors.AllowCredentials,
			MaxAge:           int(c.Cors.MaxAge),
		}))
	}
}
