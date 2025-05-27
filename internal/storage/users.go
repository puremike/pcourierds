package storage

import (
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) CreateUser(user *models.User) error {
	return nil
}
