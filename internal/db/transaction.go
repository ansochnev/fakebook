package db

import "database/sql"

type Tx interface {
	Exec(query string, args ...any) (sql.Result, error)
}
