package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type FileDetails struct {
	FileName   string
	LineNumber int
}

func (fd *FileDetails) IsValid() bool {
	return fd.FileName != "" && fd.LineNumber > 0
}

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

	for task.Title[task.HeadingNo] == '#' {
		task.HeadingNo++
	}
	for j := task.HeadingNo + 1; task.Title[j] != ':'; j++ {
		task.Type += string(task.Title[j])
	}

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
	lines, err := readLinesFromFile(task.FileDetails.FileName)
	if err != nil {
		return err
	}

	lines[task.FileDetails.LineNumber-1] = strings.Repeat("#", task.HeadingNo) + " " + task.Type + ": " + newTitle

	err = writeLinesToFile(task.FileDetails.FileName, lines)
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) UpdateProperty(propertyType PropertyType, p any) error {
	lines, err := readLinesFromFile(task.FileDetails.FileName)
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

	err = writeLinesToFile(task.FileDetails.FileName, lines)
	if err != nil {
		return err
	}
	return nil
}

func UpdateClockDataInSlice(lines []string, propertyName string, value string, hasProperties bool, taskIndex int) ([]string, error) {
	if hasProperties {
		found := false
		i := taskIndex + 1
		for i < len(lines) && !strings.Contains(lines[i], "<!-- MTASK") {
			i++
		}
		if i >= len(lines) {
			return nil, fmt.Errorf("no MTASK comment found")
		}
		startM := i

		valueToSearch := ":" + propertyName + ":"
		for ; i < len(lines) && !strings.Contains(lines[i], "-->"); i++ {
			if strings.Contains(lines[i], valueToSearch) {
				found = true
				break
			}
		}

		if i >= len(lines) {
			return nil, fmt.Errorf("closing MTASK comment not found")
		}

		if !found {
			valueToInsert := []string{fmt.Sprintf(":%s:", propertyName), value, ":END:"}
			// no property so need to add a new one
			if i == startM {
				// comment doesn't have any properties
				wrappedValue := wrapSliceInMTaskComment(valueToInsert)
				// skipping the <!-- MTASK --> line
				lines = append(lines[:i], append(wrappedValue, lines[i+1:]...)...)
			} else {
				// comment has properties, so we need to add a new line for property in the end
				lines = append(lines[:i], append(valueToInsert, lines[i:]...)...)
			}
		} else {
			// already has property, so we need to update it
			lines = append(lines[:i+1], append([]string{value}, lines[i+1:]...)...)
		}

	} else {
		valueToInsert := []string{fmt.Sprintf(":%s:", propertyName), value, ":END:"}
		sliceToAdd := wrapSliceInMTaskComment(valueToInsert)
		if len(lines) == taskIndex+1 {
			lines = append(lines, sliceToAdd...)
		} else {
			result := make([]string, len(lines)+len(sliceToAdd))
			copy(result, lines[:taskIndex+1])
			copy(result[taskIndex+1:], sliceToAdd)
			copy(result[taskIndex+1+len(sliceToAdd):], lines[taskIndex+1:])
			lines = result
		}
	}

	return lines, nil
}

func UpdatePropertyInSlice(lines []string, propertyName string, value string, hasProperties bool, taskIndex int) ([]string, error) {
	valueToInsert := fmt.Sprintf("%s: %s", propertyName, value)
	if hasProperties {
		found := false
		i := taskIndex + 1
		for i < len(lines) && !strings.Contains(lines[i], "<!-- MTASK") {
			i++
		}
		if i >= len(lines) {
			return nil, fmt.Errorf("no MTASK comment found")
		}
		startM := i

		for ; i < len(lines) && !strings.Contains(lines[i], "-->"); i++ {
			if strings.Contains(lines[i], propertyName) {
				found = true
				break
			}
		}

		if i >= len(lines) {
			return nil, fmt.Errorf("closing MTASK comment not found")
		}

		if !found {
			// no property so need to add a new one
			if i == startM {
				// comment doesn't have any properties
				wrappedValue := wrapInMTaskComment(valueToInsert)
				// skipping the <!-- MTASK --> line
				lines = append(lines[:i], append(wrappedValue, lines[i+1:]...)...)
			} else {
				// comment has properties, so we need to add a new line for property in the end
				lines = append(lines[:i], append([]string{valueToInsert}, lines[i:]...)...)
			}
		} else {
			// already has property, so we need to update it
			lines[i] = valueToInsert
		}

	} else {
		sliceToAdd := wrapInMTaskComment(valueToInsert)
		if len(lines) == taskIndex+1 {
			lines = append(lines, sliceToAdd...)
		} else {
			result := make([]string, len(lines)+len(sliceToAdd))
			copy(result, lines[:taskIndex+1])
			copy(result[taskIndex+1:], sliceToAdd)
			copy(result[taskIndex+1+len(sliceToAdd):], lines[taskIndex+1:])
			lines = result
		}
	}

	return lines, nil
}

func wrapInMTaskComment(value string) []string {
	wrappedValue := make([]string, 0, 4)
	wrappedValue = append(wrappedValue, "")
	wrappedValue = append(wrappedValue, "<!-- MTASK")
	wrappedValue = append(wrappedValue, value)
	wrappedValue = append(wrappedValue, "-->")
	return wrappedValue
}

func wrapSliceInMTaskComment(value []string) []string {
	wrappedValue := make([]string, 0, 3+len(value))
	wrappedValue = append(wrappedValue, "")
	wrappedValue = append(wrappedValue, "<!-- MTASK")
	wrappedValue = append(wrappedValue, value...)
	wrappedValue = append(wrappedValue, "-->")
	return wrappedValue
}

func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

func writeLinesToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
