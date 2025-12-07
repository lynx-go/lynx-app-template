package events

import "time"

type DemoEvent struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
