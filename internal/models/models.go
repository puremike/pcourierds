package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // "user", "dispatcher", "admin"
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Package struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	DispatcherID int       `json:"dispatcher_id"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type DispatcherApplication struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Vehicle   string    `json:"vehicle"`
	License   string    `json:"license"`
	Status    string    `json:"status"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
}
