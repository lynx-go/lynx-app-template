package main

import (
	"github.com/google/wire"
	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx-app-template/internal/api"
	"github.com/lynx-go/lynx-app-template/internal/infra"
	configpb "github.com/lynx-go/lynx-app-template/internal/pkg/config"
	"github.com/lynx-go/lynx-app-template/pkg/pubsub"
	"github.com/lynx-go/lynx/boot"
	"github.com/lynx-go/lynx/contrib/kafka"
	"github.com/lynx-go/lynx/contrib/schedule"
	"github.com/lynx-go/lynx/server/http"
)

//go:generate wire

var ProviderSet = wire.NewSet(
	boot.New,
	api.ProviderSet,
	infra.ProviderSet,
	NewComponents,
	NewComponentBuilders,
	NewOnStarts,
	NewOnStops,
	NewHealthChecks,
	NewConfig,
)

func NewConfig(app lynx.Lynx) (*configpb.AppConfig, error) {
	var c configpb.AppConfig
	if err := app.Config().Unmarshal(&c, lynx.TagNameJSON); err != nil {
		return nil, err
	}
	return &c, nil
}

func NewHealthChecks(app lynx.Lynx) lynx.HealthCheckFunc {
	return app.HealthCheckFunc()
}

func NewComponents(
	httpServer *http.Server,
	scheduler *schedule.Scheduler,
	broker *pubsub.PubSub,
	binder *kafka.Binder,
) []lynx.Component {
	return []lynx.Component{
		scheduler,
		broker,
		httpServer,
		binder,
	}
}

func NewOnStarts(
	router *pubsub.Router,
) lynx.OnStartHooks {
	hooks := lynx.OnStartHooks{
		router.Run,
	}
	return hooks
}

func NewOnStops() lynx.OnStopHooks {
	hooks := lynx.OnStopHooks{}
	return hooks
}

func NewComponentBuilders(
	binder *kafka.Binder,
) []lynx.ComponentBuilder {
	builders := []lynx.ComponentBuilder{}
	builders = append(builders, binder.Builders()...)
	return builders
}
