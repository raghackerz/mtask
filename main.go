package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type RipgrepText struct {
	Text string `json:"text"`
}

type RipgrepResult struct {
	Type string `json:"type"`
	Data struct {
		Path       RipgrepText `json:"path"`
		Lines      RipgrepText `json:"lines"`
		LineNumber int         `json:"line_number"`
	} `json:"data"`
}

func main() {
	// rg --json -g '*.md' -U -e '^#{1,6}\s+TODO:.+\n(?:\s*<!-- MTASK[\S\s]*?-->)?'
	out, err := exec.Command("rg", "--json", "-g", "*.md", "-U", "-e", "^#{1,6}\\s+TODO:.+\\n(?:\\s*<!-- MTASK[\\S\\s]*?-->)?").Output()
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
			task := Task{}
			task.FileDetails.FileName = fileName
			task.FileDetails.LineNumber = r.Data.LineNumber
			task.PopulateDetails(r.Data.Lines.Text)
			tasks = append(tasks, task)
			fmt.Println(task)
		}
	}
}
