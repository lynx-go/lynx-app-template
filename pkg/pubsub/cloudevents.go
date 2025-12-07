package pubsub

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func MustFromMessage(msg *message.Message) *cloudevents.Event {
	return lo.Must1(FromMessage(msg))
}

func FromMessage(msg *message.Message) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	if err := event.UnmarshalJSON(msg.Payload); err != nil {
		return nil, err
	}
	return &event, nil
}

func NewMessage(brokerId, name string, data any) (*message.Message, error) {
	msgId := uuid.New().String()
	event := cloudevents.NewEvent()
	event.SetID(msgId)
	event.SetType(name)
	event.SetSource(brokerId)
	event.SetTime(time.Now())
	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return nil, err
	}
	bytes, err := event.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return message.NewMessage(msgId, bytes), nil
}

func MustNewMessage(brokerId, name string, data any) *message.Message {
	return lo.Must1(NewMessage(brokerId, name, data))
}
