package db

import (
	"fmt"
	"testing"
)

type duplicateTestCase struct {
	Value  string
	Table  string
	Column string
}

func TestMustParseDuplicateError(t *testing.T) {
	// Arrange
	messageFormat := "Duplicate entry '%s' for key '%s.%s'"

	cases := []duplicateTestCase{
		{"value", "table", "column"},
		{"V", "T", "C"},
		{"значение", "таблица", "колонка"},
		{"З", "Т", "К"},
	}

	messages := []string{}
	for _, entry := range cases {
		messages = append(messages,
			fmt.Sprintf(messageFormat, entry.Value, entry.Table, entry.Column))
	}

	// Test
	for i, message := range messages {
		// Act
		columnName, value := mustParseDuplicateError(message)

		// Assert
		if columnName != cases[i].Column {
			t.Errorf("Invalid columnName: \"%s\", expect \"%s\"", columnName, cases[i].Column)
		}
		if value != cases[i].Value {
			t.Errorf("Invalid value: \"%s\", expect \"%s\"", value, cases[i].Value)
		}
	}
}
