package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

var TasksInFile map[string][]Task = make(map[string][]Task)

func main() {
	err := SyncAllTasks()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(TasksInFile)
}

func SyncAllTasks() error {
	out, err := GetAllMatchesInRootDir()
	if err != nil {
		return err
	}
	tasks, err := ParseMatchesToTasks(out)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		if _, ok := TasksInFile[task.FileDetails.FileName]; !ok {
			TasksInFile[task.FileDetails.FileName] = make([]Task, 0)
		}
		TasksInFile[task.FileDetails.FileName] = append(TasksInFile[task.FileDetails.FileName], task)
	}
	return nil
}

func SyncTasksInAFile(filepath string) error {
	out, err := GetAllMatchesInAFile(filepath)
	if err != nil {
		return err
	}
	tasks, err := ParseMatchesToTasks(out)
	if err != nil {
		return err
	}
	TasksInFile[filepath] = tasks
	return nil
}

func ParseMatchesToTasks(match []byte) ([]Task, error) {
	var fileName string
	tasks := make([]Task, 0)
	for v := range bytes.SplitSeq(match, []byte("\n")) {
		var r RipgrepResult
		if len(bytes.TrimSpace(v)) == 0 {
			continue
		}
		err1 := json.Unmarshal(v, &r)
		if err1 != nil {
			fmt.Println("Error unmarshalling JSON:", err1)
			continue
		}
		if r.Type == "begin" {
			fileName = r.Data.Path.Text
		}
		if r.Type == "match" {
			currentLineNumber := r.Data.LineNumber
			temp := 0
			for _, match := range r.Data.Submatches {
				for temp < match.Start {
					if r.Data.Lines.Text[temp] == '\n' {
						currentLineNumber++
					}
					temp++
				}
				task := Task{}
				task.FileDetails.FileName = fileName
				task.FileDetails.LineNumber = currentLineNumber
				task.PopulateDetails(match.Match.Text)
				tasks = append(tasks, task)
			}
		}
	}
	return tasks, nil
}
