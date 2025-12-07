package server

import (
	"github.com/lynx-go/lynx-template/internal/api/events"
	"github.com/lynx-go/lynx-template/pkg/pubsub"
)

func NewPubSub() *pubsub.PubSub {
	return pubsub.NewPubSub()
}

func NewPubSubRouter(
	pubSub *pubsub.PubSub,
	demo *events.DemoHandler,
) *pubsub.Router {
	return pubsub.NewRouter(pubSub, []pubsub.Handler{
		demo,
	})
}
