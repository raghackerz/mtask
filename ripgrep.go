package main

import (
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
		Submatches []struct {
			Match RipgrepText `json:"match"`
			Start int         `json:"start"`
			End   int         `json:"end"`
		} `json:"submatches"`
	} `json:"data"`
}

func GetAllMatchesInRootDir() ([]byte, error) {
	// rg --json -g '*.md' -U -e '^#{1,6}.+\n(?:\s*<!-- MTASK[\S\s]*?-->)?'
	return exec.Command("rg", "--json", "-g", "*.md", "-U", "-e", "^#{1,6}.+\\n(?:\\s*<!-- MTASK[\\S\\s]*?-->)?", RootDir).Output()
}
func GetAllMatchesInAFile(filepath string) ([]byte, error) {
	return exec.Command("rg", "--json", "-g", "*.md", "-U", "-e", "^#{1,6}.+\\n(?:\\s*<!-- MTASK[\\S\\s]*?-->)?", filepath).Output()
}
