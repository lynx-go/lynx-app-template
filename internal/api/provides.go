package api

import (
	"github.com/google/wire"
	"github.com/lynx-go/lynx-app-template/internal/api/events"
	"github.com/lynx-go/lynx-app-template/internal/api/http"
	"github.com/lynx-go/lynx-app-template/internal/api/tasks"
)

var ProviderSet = wire.NewSet(
	http.NewHelloAPI,
	tasks.NewDemoTask,
	events.NewDemoHandler,
)
