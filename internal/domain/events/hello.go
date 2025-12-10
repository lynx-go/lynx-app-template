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

type EventName string

func (e EventName) String() string {
	return string(e)
}

func (e EventName) ProducerName() string {
	return kafka.ToProducerName(e.String())
}

func (e EventName) ConsumerName() string {
	return kafka.ToConsumerName(e.String())
}

const EventNameHello EventName = "hello"
