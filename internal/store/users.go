package store

import (
	"context"
	"database/sql"

	"github.com/puremike/pcourierds/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO users (username, email, role, password) VALUES ($1, $2, $3, $4) RETURNING id, username, email, role, created_at`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Role, user.Password).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserStore) GetUserById(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, username, email, role, password, created_at FROM users WHERE id = $1`

	if err := u.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, username, email, role, password, created_at FROM users WHERE email = $1`

	if err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.Password, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserStore) UpdateUser(ctx context.Context, user *models.User, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4 RETURNING id, username, email, role, created_at`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Role, id).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt); err != nil {

	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserStore) UpdateUserRole(ctx context.Context, user *models.User, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE users SET role = $1 WHERE id = $2`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query, user.Role, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *UserStore) UpdatePassword(ctx context.Context, user *models.User, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE users SET password = $1 WHERE id = $2 RETURNING password`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, user.Password, id).Scan(&user.Password); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (u *UserStore) GetAllUsers(ctx context.Context) (*[]models.User, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `SELECT id, username, email, role, created_at FROM users`

	var users []models.User

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}
