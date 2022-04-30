package handlers

import (
	"encoding/json"
	"reflect"
	"testing"
)

func toJSON(v any) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

func TestSplitIntoParagraphs(t *testing.T) {
	testCases := []struct {
		inputText string
		expected  []string
	}{
		{"", []string{""}},
		{"Without new line", []string{"Without new line\n"}},
		{"Single line\n", []string{"Single line\n"}},
		{"Ends with multiple new lines\n\n\n", []string{"Ends with multiple new lines\n"}},
		{"Multiple\nlines\n", []string{"Multiple\nlines\n"}},
		{"Multiple\n\nParagraphs", []string{"Multiple\n", "Paragraphs\n"}},
		{"Multiple\n\n\nParagraphs\n\n", []string{"Multiple\n", "Paragraphs\n"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputText, func(t *testing.T) {
			actual := splitIntoParagraphs(testCase.inputText)
			if !reflect.DeepEqual(actual, testCase.expected) {
				t.Fatalf("\nExpect: %v,\nactual: %v",
					toJSON(testCase.expected), toJSON(actual))
			}
		})
	}
}
