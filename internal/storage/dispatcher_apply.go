package storage

import "database/sql"

type DispatcherApplyStore struct {
	db *sql.DB
}
