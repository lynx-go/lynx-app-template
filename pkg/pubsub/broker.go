package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/lynx-go/lynx/contrib/pubsub"
)

// PubSub Wrap Lynx PubSub to use cloudevents as standard event format
type PubSub struct {
	pubsub.Broker
}

func NewPubSub() *PubSub {
	broker := pubsub.NewBroker(pubsub.Options{})
	return &PubSub{broker}
}

func (b *PubSub) Publish(ctx context.Context, name string, data any, opts ...pubsub.PublishOption) error {
	o := &pubsub.PublishOptions{}
	for _, opt := range opts {
		opt(o)
	}
	msg, err := NewMessage(b.ID(), name, data)
	if err != nil {
		return err
	}

	return b.Broker.Publish(ctx, name, msg, opts...)
}

type HandlerFunc func(ctx context.Context, e *cloudevents.Event) error

func (b *PubSub) Subscribe(eventName, handlerName string, h HandlerFunc, opts ...pubsub.SubscribeOption) error {
	handler := func(ctx context.Context, msg *message.Message) error {
		event := cloudevents.NewEvent()
		if err := event.UnmarshalJSON(msg.Payload); err != nil {
			return err
		}
		return h(ctx, &event)
	}
	return b.Broker.Subscribe(eventName, handlerName, handler, opts...)
}

type Handler interface {
	EventName() string
	HandlerName() string
	HandlerFunc() HandlerFunc
}
