package pubsub

import (
	"context"

	"github.com/lynx-go/x/log"
)

type Router struct {
	handlers []Handler
	pubSub   *PubSub
}

func NewRouter(pubSub *PubSub, handlers []Handler) *Router {
	return &Router{
		pubSub:   pubSub,
		handlers: handlers,
	}
}

func (r *Router) Run(ctx context.Context) error {
	for i := range r.handlers {
		h := r.handlers[i]
		ctx := log.WithContext(ctx, "handler_name", h.HandlerName(), "event_name", h.EventName())
		log.InfoContext(ctx, "binding handler")
		if err := r.pubSub.Subscribe(h.EventName(), h.HandlerName(), h.HandlerFunc()); err != nil {
			return err
		}
	}
	return nil
}
