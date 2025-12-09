package events

import "time"

type HelloEvent struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
