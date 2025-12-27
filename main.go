package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	out, err := GetAllMatchesInRootDir()
	if err != nil {
		log.Fatal(err)
	}
	var fileName string
	tasks := make([]Task, 0)
	for v := range bytes.SplitSeq(out, []byte("\n")) {
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
				fmt.Println(task)
			}
		}
	}
}
