package main

import (
	"fmt"
	"strings"
	"time"
)

type TaskType int

const (
	Todo TaskType = iota
	Next
	Done
	NotATask
)

var TaskTypeNames = map[TaskType]string{
	Todo:     "TODO",
	Next:     "NEXT",
	Done:     "DONE",
	NotATask: "",
}
var TaskTypeValues = map[string]TaskType{
	TaskTypeNames[Todo]: Todo,
	TaskTypeNames[Next]: Next,
	TaskTypeNames[Done]: Done,
}

func (t TaskType) String() string {
	return TaskTypeNames[t]
}

func GetTaskTypeFromString(typeStr string) TaskType {
	taskType, ok := TaskTypeValues[typeStr]
	if ok {
		return taskType
	}
	return NotATask
}

type Task struct {
	Title               string
	Type                TaskType
	HeadingNo           int
	Scheduled           time.Time
	Deadline            time.Time
	ClockData           []TimeSpan
	FileDetails         FileDetails
	SubTasks            []Task
	HasPropertiesInFile bool
}

func (task *Task) String() string {
	var builder strings.Builder
	for i := 0; i < task.HeadingNo; i++ {
		builder.WriteByte('#')
	}
	builder.WriteByte(' ')
	if task.Type != NotATask {
		builder.WriteString(task.Type.String())
		builder.WriteString(": ")
	}
	builder.WriteString(task.Title)
	builder.WriteByte('\n')
	if !task.HasPropertiesInFile {
		return builder.String()
	}
	builder.WriteString("<!-- MTASK\n")
	if !task.Scheduled.IsZero() {
		builder.WriteString(fmt.Sprintf("%s: %s\n", PropertyTypeNames[Scheduled], task.Scheduled.Format(DateFormat)))
	}
	if !task.Deadline.IsZero() {
		builder.WriteString(fmt.Sprintf("%s: %s\n", PropertyTypeNames[Deadline], task.Scheduled.Format(DateFormat)))
	}
	if len(task.ClockData) > 0 {
		builder.WriteString(fmt.Sprintf(":%s:\n", PropertyTypeNames[ClockData]))
		for _, ts := range task.ClockData {
			value, _ := ts.GetFormattedValue()
			builder.WriteString(fmt.Sprintf("%s\n", value))
		}
		builder.WriteString(":END:\n")
	}
	builder.WriteString("-->\n")
	return builder.String()
}

func (task *Task) ParseProperties(propertyMap map[string]string) error {
	for _, parser := range PropertyParsers {
		value, ok := propertyMap[parser.Property]
		if ok {
			parser.Parse(task, value)
		}
	}
	return nil
}

func ParseTitleLine(str string) (int, TaskType, string) {
	headingNo := 0
	taskType := NotATask
	var titleStr string
	for str[headingNo] == '#' {
		headingNo++
	}
	indexTitleStart := headingNo + 1
	if strings.Contains(str, ":") {
		typeStr := ""
		for j := headingNo + 1; str[j] != ':'; j++ {
			typeStr += string(str[j])
		}
		indexTitleStart += len(typeStr) + 1
		typeStr = strings.TrimSpace(typeStr)
		taskType = GetTaskTypeFromString(typeStr)
	} else {
		taskType = NotATask
	}
	titleStr = strings.TrimSpace(str[indexTitleStart:])
	return headingNo, taskType, titleStr
}

func (task *Task) PopulateDetails(match string) error {
	/*if !task.FileDetails.IsValid() {
		return fmt.Errorf("invalid file details: %v", task.FileDetails)
	}
	content, err := os.ReadFile(task.FileDetails.FileName)
	if err != nil {
		return err
	}*/
	match = strings.ReplaceAll(match, "\r\n", "\n")
	task.HasPropertiesInFile = strings.Count(match, "\n") >= 2
	keyValue := make(map[string]string)
	var builder strings.Builder
	i := 0
	// read the first line, which contains the task title
	for match[i] != '\n' {
		builder.WriteByte(match[i])
		i++
	}
	task.HeadingNo, task.Type, task.Title = ParseTitleLine(builder.String())
	builder.Reset()

	if !task.HasPropertiesInFile {
		return nil
	}
	// ignore characters until we find the 'K' character
	for match[i] != 'K' {
		i++
	}
	i++

	// read the properties defined in the comment.
	for {
		// igone all spaces and new lines
		for match[i] == ' ' || match[i] == '\n' {
			i++
		}
		// check if reached end of comment, and exit
		if match[i] == '-' && match[i+1] == '-' && match[i+2] == '>' {
			break
		}

		multiline := false
		if match[i] == ':' {
			multiline = true
			// skip the first : and read till : for property name
			i++
		}
		// read the key
		for match[i] != ':' {
			if match[i] != ' ' && match[i] != '\n' {
				builder.WriteByte(match[i])
			}
			i++
		}
		key := builder.String()
		builder.Reset()
		// ignore the ':' character
		i++

		// read the value
		if multiline {
			// skip the newline
			i++
			// read till we encounder \n:END:
			for match[i] != '\n' || match[i+1] != byte(':') || match[i+2] != byte('E') || match[i+3] != byte('N') || match[i+4] != byte('D') || match[i+5] != byte(':') {
				builder.WriteByte(match[i])
				i++
			}
			i += 5
		} else {
			for match[i] != '\n' {
				if match[i] != ' ' && match[i] != '\n' {
					builder.WriteByte(match[i])
				}
				i++
			}
		}
		// ignore the new line
		i++

		keyValue[key] = builder.String()
		builder.Reset()
	}
	task.ParseProperties(keyValue)
	return nil
}

func (task *Task) WriteToFile() error {
	lines, err := ReadLinesFromFile(task.FileDetails.FileName)
	if err != nil {
		return err
	}
	startLine := task.FileDetails.LineNumber - 1
	endLine := startLine
	// TODO: Think how we make sure that this always updated with file
	if task.HasPropertiesInFile {
		startFound := true
		for endLine < len(lines) {
			if strings.Contains(lines[endLine], "<!-- MTASK") {
				startFound = true
			}
			if startFound && strings.Contains(lines[endLine], "-->") {
				break
			}
			endLine++
		}
	}
	taskString := task.String()
	taskLines := strings.Split(taskString, "\n")
	taskLines = taskLines[:len(taskLines)-1] // remove the last empty line
	resultLines := make([]string, 0)
	resultLines = append(resultLines, lines[:startLine]...)
	resultLines = append(resultLines, taskLines...)
	resultLines = append(resultLines, lines[endLine+1:]...)
	err = WriteLinesToFile(task.FileDetails.FileName+"temp.md", resultLines)
	if err != nil {
		return err
	}
	return nil
}
