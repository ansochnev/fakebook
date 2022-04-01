package backend

import "database/sql"

type Backend struct {
	db *sql.DB
}

func New(db *sql.DB) (*Backend, error) {
	return &Backend{
		db: db,
	}, nil
}
