package server

import (
	"context"

	"github.com/lynx-go/lynx-app-template/internal/api/events"
	configpb "github.com/lynx-go/lynx-app-template/internal/pkg/config"
	"github.com/lynx-go/lynx-app-template/pkg/pubsub"
	"github.com/lynx-go/lynx/contrib/kafka"
	"github.com/lynx-go/x/log"
)

func NewPubSub() *pubsub.PubSub {
	return pubsub.NewPubSub()
}

func NewPubSubRouter(
	pubSub *pubsub.PubSub,
	demo *events.HelloHandler,
) *pubsub.Router {
	return pubsub.NewRouter(pubSub, []pubsub.Handler{
		demo,
	})
}

func NewKafkaBinderForServer(pubsub *pubsub.PubSub, config *configpb.AppConfig) *kafka.Binder {
	return NewKafkaBinder(pubsub, config, false)
}

func NewKafkaBinderForCli(pubsub *pubsub.PubSub, config *configpb.AppConfig) *kafka.Binder {
	return NewKafkaBinder(pubsub, config, true)
}

func NewKafkaBinder(pubsub *pubsub.PubSub, config *configpb.AppConfig, disableSub bool) *kafka.Binder {
	bindOptions := kafka.BinderOptions{
		SubscribeOptions: map[string]kafka.ConsumerOptions{},
		PublishOptions:   map[string]kafka.ProducerOptions{},
	}
	if config.Pubsub != nil && config.Pubsub.Kafka != nil {
		cfgs := config.Pubsub.Kafka
		for k, c := range cfgs {
			if c.Consumer != nil && !disableSub {
				bindOptions.SubscribeOptions[k] = kafka.ConsumerOptions{
					Brokers:          c.Brokers,
					Topic:            c.Topic,
					Group:            c.Consumer.GroupId,
					ErrorHandlerFunc: logError,
					Instances:        int(c.Consumer.Instances),
					LogMessage:       c.Consumer.LogMessage,
				}
			}
			if c.Producer != nil {
				bindOptions.PublishOptions[k] = kafka.ProducerOptions{
					Brokers:    c.Brokers,
					Topic:      c.Topic,
					LogMessage: c.Producer.LogMessage,
				}
			}
		}
	}

	return kafka.NewBinder(bindOptions, pubsub.Broker)
}

func logError(err error) error {
	log.ErrorContext(context.TODO(), "handle kafka pubsub error", err)
	return nil
}
