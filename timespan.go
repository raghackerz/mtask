package main

import (
	"time"
)

type TimeSpan struct {
	StartTime time.Time
	EndTime   time.Time
	IsRunning bool
}
