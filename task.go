package main

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Title         string
	Type          string
	HeadingNo     int
	Scheduled     time.Time
	Deadline      time.Time
	ClockData     []TimeSpan
	FileDetails   FileDetails
	SubTasks      []Task
	HasProperties bool
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

func (task *Task) PopulateDetails(match string) error {
	/*if !task.FileDetails.IsValid() {
		return fmt.Errorf("invalid file details: %v", task.FileDetails)
	}
	content, err := os.ReadFile(task.FileDetails.FileName)
	if err != nil {
		return err
	}*/
	task.HasProperties = strings.Count(match, "\n") >= 2
	keyValue := make(map[string]string)
	var builder strings.Builder
	i := 0
	// read the first line, which contains the task title
	for match[i] != '\n' {
		builder.WriteByte(match[i])
		i++
	}
	task.Title = builder.String()
	builder.Reset()

	// check the level of heading
	task.HeadingNo = 0
	for task.Title[task.HeadingNo] == '#' {
		task.HeadingNo++
	}

	// type of task
	for j := task.HeadingNo + 1; task.Title[j] != ':'; j++ {
		task.Type += string(task.Title[j])
	}
	task.Type = strings.TrimSpace(task.Type)

	if !task.HasProperties {
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

func (task *Task) ChangeTitle(newTitle string) error {
	lines, err := ReadLinesFromFile(task.FileDetails.FileName)
	if err != nil {
		return err
	}

	lines[task.FileDetails.LineNumber-1] = strings.Repeat("#", task.HeadingNo) + " " + task.Type + ": " + newTitle

	err = WriteLinesToFile(task.FileDetails.FileName, lines)
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) UpdateProperty(propertyType PropertyType, p any) error {
	lines, err := ReadLinesFromFile(task.FileDetails.FileName)
	if err != nil {
		return err
	}

	var value string
	if propertyType == Scheduled || propertyType == Deadline {
		newTime, ok := p.(time.Time)
		if !ok {
			return fmt.Errorf("expected time.Time, got %T", p)
		}
		value = newTime.Format(DateFormat)
	}

	lines, err = UpdatePropertyInSlice(lines, PropertyTypeNames[propertyType], value, task.HasProperties, task.FileDetails.LineNumber-1)
	if err != nil {
		return err
	}

	err = WriteLinesToFile(task.FileDetails.FileName, lines)
	if err != nil {
		return err
	}
	return nil
}
