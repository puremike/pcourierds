package store

import (
	"context"
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type DispatcherApplyStore struct {
	db *sql.DB
}

func (d *DispatcherApplyStore) DispatcherApplication(ctx context.Context, application *models.DispatcherApplication) (*models.DispatcherApplication, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO dispatchers_apply (user_id, vehicle_type, vehicle_plate_number, vehicle_year, vehicle_model, driver_license, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, vehicle_type, vehicle_plate_number, vehicle_year, vehicle_model, driver_license, status, created_at`

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, application.UserID, application.VehicleType, application.VehiclePlateNumber, application.VehicleYear, application.VehicleModel, application.DriverLicense, application.Status).Scan(&application.ID, &application.UserID, &application.VehicleType, &application.VehiclePlateNumber, &application.VehicleYear, &application.VehicleModel, &application.DriverLicense, &application.Status, &application.CreatedAt); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return application, nil
}

func (d *DispatcherApplyStore) GetAllApplications(ctx context.Context) (*[]models.DispatcherApplication, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `SELECT id, user_id, vehicle_type, vehicle_plate_number, vehicle_year, vehicle_model, driver_license, status, created_at FROM dispatchers_apply`

	var dispatchersApp []models.DispatcherApplication

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d models.DispatcherApplication
		if err = rows.Scan(&d.ID, &d.UserID, &d.VehicleType, &d.VehiclePlateNumber, &d.VehicleYear, &d.VehicleModel, &d.DriverLicense, &d.Status, &d.CreatedAt); err != nil {
			return nil, err
		}

		dispatchersApp = append(dispatchersApp, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &dispatchersApp, nil
}

func (d *DispatcherApplyStore) GetApplicationById(ctx context.Context, id string) (*models.DispatcherApplication, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	dispatcherApp := &models.DispatcherApplication{}

	query := `SELECT id, user_id, vehicle_type, vehicle_plate_number, vehicle_year, vehicle_model, driver_license, status, created_at FROM dispatchers_apply WHERE id = $1`

	if err := d.db.QueryRowContext(ctx, query, id).Scan(&dispatcherApp.ID, &dispatcherApp.UserID, &dispatcherApp.VehicleType, &dispatcherApp.VehiclePlateNumber, &dispatcherApp.VehicleYear, &dispatcherApp.VehicleModel, &dispatcherApp.DriverLicense, &dispatcherApp.Status, &dispatcherApp.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDispatcherApplicationNotFound
		}
		return nil, err
	}

	return dispatcherApp, nil
}

func (d *DispatcherApplyStore) GetApplicationByUserId(ctx context.Context, userId string) (*models.DispatcherApplication, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	dispatcherApp := &models.DispatcherApplication{}

	query := `SELECT id, user_id, vehicle_type, vehicle_plate_number, vehicle_year, vehicle_model, driver_license, status, created_at FROM dispatchers_apply WHERE user_id = $1`

	if err := d.db.QueryRowContext(ctx, query, userId).Scan(&dispatcherApp.ID, &dispatcherApp.UserID, &dispatcherApp.VehicleType, &dispatcherApp.VehiclePlateNumber, &dispatcherApp.VehicleYear, &dispatcherApp.VehicleModel, &dispatcherApp.DriverLicense, &dispatcherApp.Status, &dispatcherApp.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDispatcherApplicationNotFound
		}
		return nil, err
	}

	return dispatcherApp, nil
}

func (d *DispatcherApplyStore) DeleteApplicationByUserId(ctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `DELETE FROM dispatchers_apply WHERE user_id = $1`

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, userId); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *DispatcherApplyStore) UpdateDispatchApplicationStatus(ctx context.Context, dispatch *models.DispatcherApplication, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE dispatchers_apply SET status = $1 WHERE id = $2`

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query, dispatch.Status, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
