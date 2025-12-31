package main

import (
	"strings"
	"time"
)

func ParseScheduleFromString(value string) (time.Time, error) {
	return time.ParseInLocation(DateFormat, value, time.Local)
}

func TryParseSchedule(task *Task, value string) error {
	val, err := ParseScheduleFromString(value)
	if err != nil {
		return err
	}
	task.Scheduled = val
	return nil
}

func ParseDeadlineFromString(value string) (time.Time, error) {
	return time.ParseInLocation(DateFormat, value, time.Local)
}
func TryParseDeadline(task *Task, value string) error {
	val, err := ParseDeadlineFromString(value)
	if err != nil {
		return err
	}
	task.Deadline = val
	return nil
}

func ParseClockDataFromString(value string) ([]TimeSpan, error) {
	res := make([]TimeSpan, 0)
	for timeSpan := range strings.SplitSeq(value, "\n") {
		times := strings.Split(timeSpan, "-")
		for i := range times {
			times[i] = strings.TrimSpace(times[i])
		}
		startTime, err := time.ParseInLocation("2006/01/02 15:04", times[0], time.Local)
		if err != nil {
			return nil, err
		}

		if len(times) < 2 || times[1] == "" {
			res = append(res, TimeSpan{StartTime: startTime, IsRunning: true})
		} else {
			endTime, err := time.ParseInLocation("2006/01/02 15:04", times[1], time.Local)
			if err != nil {
				return nil, err
			}
			res = append(res, TimeSpan{StartTime: startTime, EndTime: endTime, IsRunning: false})
		}
	}
	return res, nil
}

func TryParseClockData(task *Task, value string) error {
	parsed, err := ParseClockDataFromString(value)
	if err != nil {
		return err
	}
	task.ClockData = parsed
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
