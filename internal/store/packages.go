package store

import (
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type PackageStore struct {
	db *sql.DB
}

func (p *PackageStore) CreatePackage(pack *models.Package) error {
	return nil
}
