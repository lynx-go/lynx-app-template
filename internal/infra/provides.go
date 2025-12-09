package infra

import (
	"github.com/google/wire"
	"github.com/lynx-go/lynx-app-template/internal/infra/server"
)

var ProviderSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewRouter,
	server.NewScheduler,
	server.NewPubSub,
	server.NewPubSubRouter,
)
