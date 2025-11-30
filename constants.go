package main

const DateFormat = "2006/01/02"
const TimeFormat = "2006/01/02 15:04"

const rootDir = "."

type PropertyType int

const (
	Scheduled PropertyType = iota
	Deadline
	ClockData
)

var PropertyTypeNames = map[PropertyType]string{
	Scheduled: "SCHEDULED",
	Deadline:  "DEADLINE",
	ClockData: "CLOCK_DATA",
}
