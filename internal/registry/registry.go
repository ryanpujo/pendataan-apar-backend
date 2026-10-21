package registry

import (
	"database/sql"
)

type registry struct {
	db *sql.DB
}

func NewRegistry(db *sql.DB) *registry {
	return &registry{
		db: db,
	}
}
