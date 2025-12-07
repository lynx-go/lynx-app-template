package events

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/lynx-go/lynx-template/internal/domain/events"
	"github.com/lynx-go/lynx-template/pkg/pubsub"
	"github.com/lynx-go/x/encoding/json"
	"github.com/lynx-go/x/log"
)

type DemoHandler struct {
}

func (h *DemoHandler) EventName() string {
	return "demo"
}

func (h *DemoHandler) HandlerName() string {
	return "DemoHandler"
}

func (h *DemoHandler) HandlerFunc() pubsub.HandlerFunc {
	return func(ctx context.Context, e *cloudevents.Event) error {
		data := &events.DemoEvent{}
		if err := e.DataAs(data); err != nil {
			return err
		}
		log.InfoContext(ctx, "recv demo event", "event_data", json.MustMarshalToString(data))
		return nil
	}
}

var _ pubsub.Handler = new(DemoHandler)

func NewDemoHandler() *DemoHandler {
	return &DemoHandler{}
}
