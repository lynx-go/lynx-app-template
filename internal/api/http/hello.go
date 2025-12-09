package http

import (
	"context"
	"time"

	"github.com/lynx-go/lynx-app-template/internal/domain/events"
	"github.com/lynx-go/lynx-app-template/pkg/pubsub"
	"github.com/lynx-go/lynx/contrib/kafka"
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
	if err := api.pubsub.Publish(ctx, kafka.ToProducerName("hello"), &events.HelloEvent{
		Message: req.Message,
		Time:    time.Now(),
	}); err != nil {
		log.ErrorContext(ctx, "failed to publish hello event", err)
	}
	return &HelloResult{Message: "Hello " + req.Message}, nil
}
