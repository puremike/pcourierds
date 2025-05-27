package storage

import "database/sql"

type PackageStore struct {
	db *sql.DB
}
