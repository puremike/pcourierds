package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/puremike/pcourierds/internal/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
}

type DispatchersApplyRepository interface {
	CreateDispatcherApplication(application *models.DispatcherApplication) error
}

type DispatchersRepository interface {
	CreateDispatcher(dispatcher *models.Dispatcher) error
}

type PackagesRepository interface {
	CreatePackage(pack *models.Package) error
}

type Storage struct {
	Users                  UsersRepository
	DispatcherApplications DispatchersApplyRepository
	Dispatchers            DispatchersRepository
	Packages               PackagesRepository
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:                  &UserStore{db},
		DispatcherApplications: &DispatcherApplyStore{db},
		Dispatchers:            &DispatcherStore{db},
		Packages:               &PackageStore{db},
	}
}

var (
	QueryBackgroundTimeout = 5 * time.Second
)
