package main

import (
	"fmt"
	"time"
)

type TimeSpan struct {
	StartTime time.Time
	EndTime   time.Time
	IsRunning bool
}

func (ts *TimeSpan) IsValid() bool {
	if ts.StartTime.IsZero() {
		return false
	}
	if ts.EndTime.IsZero() && !ts.IsRunning {
		return false
	}
	if ts.IsRunning && !ts.EndTime.IsZero() {
		return false
	}
	if !ts.EndTime.IsZero() && ts.EndTime.Before(ts.StartTime) {
		return false
	}
	return true
}

func (ts *TimeSpan) Format() (string, error) {
	if !ts.IsValid() {
		return "", fmt.Errorf("invalid TimeSpan: %v", ts)
	}
	if ts.IsRunning {
		return fmt.Sprintf("%s-", ts.StartTime.Format(TimeFormat)), nil
	}
	return fmt.Sprintf("%s-%s", ts.StartTime.Format(TimeFormat), ts.EndTime.Format(TimeFormat)), nil
}
