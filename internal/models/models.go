package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // "user", "dispatcher", "admin"
	Password  string    `json:"_"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Package struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	DispatcherID string    `json:"dispatcher_id"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DispatcherApplication struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	VehicleType        string    `json:"vehicle_type"`
	VehiclePlateNumber string    `json:"vehicle_plate_number"`
	VehicleYear        int       `json:"vehicle_year"`
	VehicleModel       string    `json:"vehicle_model"`
	DriverLicense      string    `json:"driver_license"`
	Status             string    `json:"status"` // pending, approved, rejected
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Dispatcher struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	ApplicationID      string    `json:"application_id"`
	VehicleType        string    `json:"vehicle_type"`
	VehiclePlateNumber string    `json:"vehicle_plate_number"`
	VehicleYear        int       `json:"vehicle_year"`
	VehicleModel       string    `json:"vehicle_model"`
	DriverLicense      string    `json:"driver_license"`
	ApprovedAt         time.Time `json:"approved_at"`
	IsActive           bool      `json:"is_active"` // Indicates if currently working
	Rating             float32   `json:"rating"`    // optional
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
