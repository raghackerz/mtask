package main

import (
	"strings"
	"time"
)

func TryParseSchedule(task *Task, value string) error {
	val, err := time.ParseInLocation("2006/01/02", value, time.Local)
	if err != nil {
		return err
	}
	task.Scheduled = val
	return nil
}

func TryParseDeadline(task *Task, value string) error {
	val, err := time.ParseInLocation("2006/01/02", value, time.Local)
	if err != nil {
		return err
	}
	task.Deadline = val
	return nil
}

func TryParseClockData(task *Task, value string) error {
	for timeSpan := range strings.SplitSeq(value, "\n") {
		times := strings.Split(timeSpan, "-")
		for i := range times {
			times[i] = strings.TrimSpace(times[i])
		}
		startTime, err := time.ParseInLocation("2006/01/02 15:04", times[0], time.Local)
		if err != nil {
			return err
		}
		endTime, err := time.ParseInLocation("2006/01/02 15:04", times[1], time.Local)
		if err != nil {
			return err
		}
		task.ClockData = append(task.ClockData, TimeSpan{StartTime: startTime, EndTime: endTime})
	}
	return nil
}

type TaskPropertyParse struct {
	Parse    func(task *Task, value string) error
	Property string
}

var PropertyParsers = []TaskPropertyParse{
	{TryParseSchedule, PropertyTypeNames[Scheduled]},
	{TryParseDeadline, PropertyTypeNames[Deadline]},
	{TryParseClockData, PropertyTypeNames[ClockData]},
}
