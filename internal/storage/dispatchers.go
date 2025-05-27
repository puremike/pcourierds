package storage

import "database/sql"

type DispatcherStore struct {
	db *sql.DB
}
