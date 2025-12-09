package infra

import (
	"github.com/google/wire"
	"github.com/lynx-go/lynx-app-template/internal/infra/clients"
	"github.com/lynx-go/lynx-app-template/internal/infra/repoimpl"
	"github.com/lynx-go/lynx-app-template/internal/infra/server"
)

var ProviderSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewRouter,
	server.NewScheduler,
	server.NewPubSub,
	server.NewPubSubRouter,
	server.NewKafkaBinderForServer,
	clients.NewDataClients,
	repoimpl.NewUserRepo,
	repoimpl.NewRuntimeVars,
)
