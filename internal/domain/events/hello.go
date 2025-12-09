package events

import (
	"time"

	"github.com/lynx-go/lynx/contrib/kafka"
)

type HelloEvent struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

const EventNameHello = "hello"

var (
	ProducerNameHello = kafka.ToProducerName(EventNameHello)
	ConsumerNameHello = kafka.ToConsumerName(EventNameHello)
)
