package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/puremike/pcourierds/internal/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User, id string) (*models.User, error)
	UpdatePassword(ctx context.Context, user *models.User, id string) error
	GetAllUsers(ctx context.Context) (*[]models.User, error)
}

type DispatchersApplyRepository interface {
	DispatcherApplication(ctx context.Context, application *models.DispatcherApplication) (*models.DispatcherApplication, error)
	GetAllApplications(ctx context.Context) (*[]models.DispatcherApplication, error)
	GetApplicationById(ctx context.Context, id string) (*models.DispatcherApplication, error)
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
	QueryBackgroundTimeout           = 5 * time.Second
	ErrUserNotFound                  = errors.New("user not found")
	ErrDispatcherApplicationNotFound = errors.New("dispatcher application not found")
)
