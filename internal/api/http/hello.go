package http

import (
	"context"
	"errors"
	"time"

	"github.com/lynx-go/lynx-app-template/internal/domain/events"
	"github.com/lynx-go/lynx-app-template/pkg/pubsub"
	"github.com/lynx-go/lynx-app-template/pkg/session"
	"github.com/lynx-go/x/log"
)

type HelloAPI struct {
	pubsub *pubsub.PubSub
}

func NewHelloAPI(
	pubsub *pubsub.PubSub,
) *HelloAPI {
	return &HelloAPI{pubsub: pubsub}
}

type HelloRequest struct {
	Message string `json:"message"`
}

type HelloResult struct {
	Message string `json:"message"`
}

func (api *HelloAPI) Hello(ctx context.Context, req *HelloRequest) (*HelloResult, error) {
	userId := session.CurrentUser(ctx)
	if userId == 0 {
		return nil, errors.New("not login")
	}
	if err := api.pubsub.Publish(ctx, events.EventNameHello, "hello", &events.HelloEvent{
		User:    userId,
		Message: req.Message,
		Time:    time.Now(),
	}); err != nil {
		log.ErrorContext(ctx, "failed to publish hello event", err)
	}
	return &HelloResult{Message: "Hello " + req.Message}, nil
}
