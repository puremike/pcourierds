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

	query := `INSERT INTO dispatchers_apply (user_id, vehicle, license, status) VALUES ($1, $2, $3, $4) RETURNING id, user_id, vehicle, license, status, created_at`

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, application.UserID, application.Vehicle, application.License, application.Status).Scan(&application.ID, &application.UserID, &application.Vehicle, &application.License, &application.Status, &application.CreatedAt); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return application, nil
}
