package main

import (
	"reflect"
	"testing"
)

// For deep comparison of slices

func TestUpdateSlice(t *testing.T) {
	tests := []struct {
		name          string
		propertyName  string
		value         string
		hasProperties bool
		taskIndex     int
		input         []string
		expected      []string
	}{
		{
			name:          "",
			propertyName:  "SCHEDULED",
			value:         "2025/07/20",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK", "SCHEDULED: 2025/07/19", "-->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "SCHEDULED: 2025/07/20", "-->"},
		},
		{
			name:          "CommentButCurrentPropertyNotPresent",
			propertyName:  "SCHEDULED",
			value:         "2025/07/20",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK", "DEADLINE: 2025/07/19", "-->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "DEADLINE: 2025/07/19", "SCHEDULED: 2025/07/20", "-->"},
		},
		{
			name:          "MultilineCommentButNoProperties",
			propertyName:  "SCHEDULED",
			value:         "2025/07/20",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK", "-->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "SCHEDULED: 2025/07/20", "-->"},
		},
		{
			name:          "CommentButNoProperties",
			propertyName:  "SCHEDULED",
			value:         "2025/07/20",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "<!-- MTASK -->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "SCHEDULED: 2025/07/20", "-->"},
		},
		{
			name:          "NoPropertiesTaskInLastLine",
			propertyName:  "SCHEDULED",
			value:         "2025/07/20",
			hasProperties: false,
			taskIndex:     0,
			input:         []string{"# TODO: testing"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "SCHEDULED: 2025/07/20", "-->"},
		},
		{
			name:          "NoProperties",
			propertyName:  "SCHEDULED",
			value:         "2025/07/20",
			hasProperties: false,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "some random text"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "SCHEDULED: 2025/07/20", "-->", "", "some random text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([]string, len(tt.input))
			copy(actual, tt.input)

			updatedLines, err := UpdatePropertyInSlice(actual, tt.propertyName, tt.value, tt.hasProperties, tt.taskIndex) // Call the function under test
			if err != nil {
				t.Errorf("UpdatePropertyInSlice failed: %v", err)
			}

			// Compare the actual modified slice with the expected slice
			if !reflect.DeepEqual(updatedLines, tt.expected) {
				t.Errorf("For input %v, expected %v, but got %v", tt.input, tt.expected, updatedLines)
			}
		})
	}
}

func TestUpdateClockDataInSlice(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		hasProperties bool
		taskIndex     int
		input         []string
		expected      []string
	}{
		{
			name:          "",
			value:         "2025/07/20-2025/07/21",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK", ":CLOCK_DATA:", "2025/07/19-2025/07/19", ":END:", "-->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", ":CLOCK_DATA:", "2025/07/20-2025/07/21", "2025/07/19-2025/07/19", ":END:", "-->"},
		},
		{
			name:          "CommentButCurrentPropertyNoPresent",
			value:         "2025/07/20-2025/07/21",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK", "DEADLINE: 2025/07/19", "-->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", "DEADLINE: 2025/07/19", ":CLOCK_DATA:", "2025/07/20-2025/07/21", ":END:", "-->"},
		},
		{
			name:          "MultilineCommentButNoProperties",
			value:         "2025/07/20-2025/07/21",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK", "-->"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", ":CLOCK_DATA:", "2025/07/20-2025/07/21", ":END:", "-->"},
		},
		{
			name:          "CommentButNoProperties",
			value:         "2025/07/20-2025/07/21",
			hasProperties: true,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "", "<!-- MTASK -->"},
			// TODO: check this if needs to be fixed as we are adding additinal blank line
			expected: []string{"# TODO: testing", "", "", "<!-- MTASK", ":CLOCK_DATA:", "2025/07/20-2025/07/21", ":END:", "-->"},
		},
		{
			name:          "NoPropertiesTaskInLastLine",
			value:         "2025/07/20-2025/07/21",
			hasProperties: false,
			taskIndex:     0,
			input:         []string{"# TODO: testing"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", ":CLOCK_DATA:", "2025/07/20-2025/07/21", ":END:", "-->"},
		},
		{
			name:          "NoProperties",
			value:         "2025/07/20-2025/07/21",
			hasProperties: false,
			taskIndex:     0,
			input:         []string{"# TODO: testing", "some random text"},
			expected:      []string{"# TODO: testing", "", "<!-- MTASK", ":CLOCK_DATA:", "2025/07/20-2025/07/21", ":END:", "-->", "some random text"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([]string, len(tt.input))
			copy(actual, tt.input)

			updatedLines, err := UpdateClockDataInSlice(actual, "CLOCK_DATA", tt.value, tt.hasProperties, tt.taskIndex) // Call the function under test
			if err != nil {
				t.Errorf("UpdateClockDataInSlice failed: %v", err)
			}

			// Compare the actual modified slice with the expected slice
			if !reflect.DeepEqual(updatedLines, tt.expected) {
				t.Errorf("For input %v, expected %v, but got %v", tt.input, tt.expected, updatedLines)
			}
		})
	}
}
