package events

import (
	"time"

	"github.com/lynx-go/lynx-app-template/pkg/bigid"
	"github.com/lynx-go/lynx/contrib/kafka"
)

type HelloEvent struct {
	User    bigid.ID  `json:"user"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

const EventNameHello = "hello"

var (
	ProducerNameHello = kafka.ToProducerName(EventNameHello)
	ConsumerNameHello = kafka.ToConsumerName(EventNameHello)
)
