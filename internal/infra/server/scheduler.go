package server

import (
	"github.com/lynx-go/lynx/contrib/schedule"
)

func NewScheduler() (*schedule.Scheduler, error) {
	return schedule.NewScheduler([]schedule.Task{}, schedule.WithDebugEnabled())
}
