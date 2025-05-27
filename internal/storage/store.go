package storage

import "database/sql"

type Storage struct {
	Users                  interface{}
	DispatcherApplications interface{}
	Dispatchers            interface{}
	Packages               interface{}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:                  UserStore{db},
		DispatcherApplications: DispatcherApplyStore{db},
		Dispatchers:            DispatcherStore{db},
		Packages:               PackageStore{db},
	}
}
