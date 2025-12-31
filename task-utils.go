package main

// import (
// 	"fmt"
// 	"strings"
// )
//
// func (task *Task) ChangeTitle(newTitle string) error {
// 	lines, err := ReadLinesFromFile(task.FileDetails.FileName)
// 	if err != nil {
// 		return err
// 	}
//
// 	if task.Type != NotATask {
// 		lines[task.FileDetails.LineNumber-1] = strings.Repeat("#", task.HeadingNo) + " " + task.Type.String() + ": " + newTitle
// 	} else {
// 		lines[task.FileDetails.LineNumber-1] = strings.Repeat("#", task.HeadingNo) + " " + newTitle
// 	}
//
// 	err = WriteLinesToFile(task.FileDetails.FileName, lines)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (task *Task) UpdateProperty(propertyType PropertyType, p any) error {
// 	lines, err := ReadLinesFromFile(task.FileDetails.FileName)
// 	if err != nil {
// 		return err
// 	}
//
// 	var value string
// 	if propertyType == Scheduled || propertyType == Deadline {
// 		newTime, ok := p.(time.Time)
// 		if !ok {
// 			return fmt.Errorf("expected time.Time, got %T", p)
// 		}
// 		value = newTime.Format(DateFormat)
// 	}
//
// 	lines, err = UpdatePropertyInSlice(lines, PropertyTypeNames[propertyType], value, task.HasPropertiesInFile, task.FileDetails.LineNumber-1)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = WriteLinesToFile(task.FileDetails.FileName+".temp", lines)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// // TODO: test this function
// func (task *Task) UpdateClockData(ts TimeSpan) error {
// 	lines, err := ReadLinesFromFile(task.FileDetails.FileName)
// 	if err != nil {
// 		return err
// 	}
//
// 	value, err := ts.GetFormattedValue()
// 	if err != nil {
// 		return err
// 	}
//
// 	lines, err = UpdateClockDataInSlice(lines, PropertyTypeNames[ClockData], value, task.HasPropertiesInFile, task.FileDetails.LineNumber-1)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = WriteLinesToFile(task.FileDetails.FileName, lines)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func UpdateClockDataInSlice(lines []string, propertyName string, value string, hasProperties bool, taskIndex int) ([]string, error) {
// 	if hasProperties {
// 		found := false
// 		i := taskIndex + 1
// 		for i < len(lines) && !strings.Contains(lines[i], "<!-- MTASK") {
// 			i++
// 		}
// 		if i >= len(lines) {
// 			return nil, fmt.Errorf("no MTASK comment found")
// 		}
// 		startM := i
//
// 		valueToSearch := ":" + propertyName + ":"
// 		for ; i < len(lines) && !strings.Contains(lines[i], "-->"); i++ {
// 			if strings.Contains(lines[i], valueToSearch) {
// 				found = true
// 				break
// 			}
// 		}
//
// 		if i >= len(lines) {
// 			return nil, fmt.Errorf("closing MTASK comment not found")
// 		}
//
// 		if !found {
// 			valueToInsert := []string{fmt.Sprintf(":%s:", propertyName), value, ":END:"}
// 			// no property so need to add a new one
// 			if i == startM {
// 				// comment doesn't have any properties
// 				wrappedValue := wrapSliceInMTaskComment(valueToInsert)
// 				// skipping the <!-- MTASK --> line
// 				lines = append(lines[:i], append(wrappedValue, lines[i+1:]...)...)
// 			} else {
// 				// comment has properties, so we need to add a new line for property in the end
// 				lines = append(lines[:i], append(valueToInsert, lines[i:]...)...)
// 			}
// 		} else {
// 			// already has property, so we need to update it
// 			lines = append(lines[:i+1], append([]string{value}, lines[i+1:]...)...)
// 		}
//
// 	} else {
// 		valueToInsert := []string{fmt.Sprintf(":%s:", propertyName), value, ":END:"}
// 		sliceToAdd := wrapSliceInMTaskComment(valueToInsert)
// 		if len(lines) == taskIndex+1 {
// 			lines = append(lines, sliceToAdd...)
// 		} else {
// 			result := make([]string, len(lines)+len(sliceToAdd))
// 			copy(result, lines[:taskIndex+1])
// 			copy(result[taskIndex+1:], sliceToAdd)
// 			copy(result[taskIndex+1+len(sliceToAdd):], lines[taskIndex+1:])
// 			lines = result
// 		}
// 	}
//
// 	return lines, nil
// }
//
// func UpdatePropertyInSlice(lines []string, propertyName string, value string, hasProperties bool, taskIndex int) ([]string, error) {
// 	valueToInsert := fmt.Sprintf("%s: %s", propertyName, value)
// 	if hasProperties {
// 		found := false
// 		i := taskIndex + 1
// 		for i < len(lines) && !strings.Contains(lines[i], "<!-- MTASK") {
// 			i++
// 		}
// 		if i >= len(lines) {
// 			return nil, fmt.Errorf("no MTASK comment found")
// 		}
// 		startM := i
//
// 		for ; i < len(lines) && !strings.Contains(lines[i], "-->"); i++ {
// 			if strings.Contains(lines[i], propertyName) {
// 				found = true
// 				break
// 			}
// 		}
//
// 		if i >= len(lines) {
// 			return nil, fmt.Errorf("closing MTASK comment not found")
// 		}
//
// 		if !found {
// 			// no property so need to add a new one
// 			if i == startM {
// 				// comment doesn't have any properties
// 				wrappedValue := wrapInMTaskComment(valueToInsert)
// 				// skipping the <!-- MTASK --> line
// 				lines = append(lines[:i], append(wrappedValue, lines[i+1:]...)...)
// 			} else {
// 				// comment has properties, so we need to add a new line for property in the end
// 				lines = append(lines[:i], append([]string{valueToInsert}, lines[i:]...)...)
// 			}
// 		} else {
// 			// already has property, so we need to update it
// 			lines[i] = valueToInsert
// 		}
//
// 	} else {
// 		sliceToAdd := wrapInMTaskComment(valueToInsert)
// 		if len(lines) == taskIndex+1 {
// 			lines = append(lines, sliceToAdd...)
// 		} else {
// 			result := make([]string, len(lines)+len(sliceToAdd))
// 			copy(result, lines[:taskIndex+1])
// 			copy(result[taskIndex+1:], sliceToAdd)
// 			copy(result[taskIndex+1+len(sliceToAdd):], lines[taskIndex+1:])
// 			lines = result
// 		}
// 	}
//
// 	return lines, nil
// }
//
// func wrapInMTaskComment(value string) []string {
// 	wrappedValue := make([]string, 0, 4)
// 	wrappedValue = append(wrappedValue, "")
// 	wrappedValue = append(wrappedValue, "<!-- MTASK")
// 	wrappedValue = append(wrappedValue, value)
// 	wrappedValue = append(wrappedValue, "-->")
// 	return wrappedValue
// }
//
// func wrapSliceInMTaskComment(value []string) []string {
// 	wrappedValue := make([]string, 0, 3+len(value))
// 	wrappedValue = append(wrappedValue, "")
// 	wrappedValue = append(wrappedValue, "<!-- MTASK")
// 	wrappedValue = append(wrappedValue, value...)
// 	wrappedValue = append(wrappedValue, "-->")
// 	return wrappedValue
// }
