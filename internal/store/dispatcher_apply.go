package store

import (
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type DispatcherApplyStore struct {
	db *sql.DB
}

func (d *DispatcherApplyStore) CreateDispatcherApplication(application *models.DispatcherApplication) error {
	return nil
}
