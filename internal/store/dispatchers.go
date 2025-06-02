package store

import (
	"context"
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type DispatcherStore struct {
	db *sql.DB
}

func (dp *DispatcherStore) CreateDispatcher(ctx context.Context, dispatcher *models.Dispatcher) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO dispatchers (user_id, application_id, vehicle_type, vehicle_plate_number, vehicle_year, vehicle_model, driver_license, approved_at, isactive, rating) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	tx, err := dp.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query,
		dispatcher.UserID,
		dispatcher.ApplicationID,
		dispatcher.VehicleType,
		dispatcher.VehiclePlateNumber,
		dispatcher.VehicleYear,
		dispatcher.VehicleModel,
		dispatcher.DriverLicense,
		dispatcher.ApprovedAt,
		dispatcher.IsActive,
		dispatcher.Rating,
	); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
