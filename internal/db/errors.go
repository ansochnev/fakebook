package db

import (
	"bytes"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const (
	_ER_DUP_ENTRY = 1062
)

func IsDuplicateError(err error) bool {
	e, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return e.Number == _ER_DUP_ENTRY
}

func MustParseDuplicateError(err error) (columnName, value string) {
	e, _ := err.(*mysql.MySQLError)
	return mustParseDuplicateError(e.Message)
}

func mustParseDuplicateError(message string) (columnName, value string) {
	// Format: Duplicate entry 'value' for key 'table.column'
	const valueStartPos = 17
	b := []byte(message[valueStartPos:])

	pos := bytes.IndexByte(b, '\'')
	if pos == -1 {
		panic(unexpectedMessageError(message))
	}

	value = string(b[:pos])
	b = b[pos+1:]

	pos = bytes.IndexByte(b, '.')
	if pos == -1 {
		panic(unexpectedMessageError(message))
	}

	n := len(b)
	columnName = string(b[pos+1 : n-1])

	return columnName, value
}

func unexpectedMessageError(message string) error {
	return fmt.Errorf("unexpected error message from DB: %v", message)
}
