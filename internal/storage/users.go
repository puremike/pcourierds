package storage

import "database/sql"

type UserStore struct {
	db *sql.DB
}
