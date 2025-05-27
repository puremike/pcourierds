package store

import (
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type DispatcherStore struct {
	db *sql.DB
}

func (d *DispatcherStore) CreateDispatcher(dispatcher *models.Dispatcher) error {
	return nil
}
