package database

import "time"

type Role string

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	RoleUser      Role = "USER"
	RoleAuthority Role = "AUTHORITY"
)

type User struct {
	ID        string    `json:"id"         db:"id"`
	Name      string    `json:"name"       db:"name"`
	Age       int       `json:"age"        db:"age"`
	Language  string    `json:"language"   db:"language"`
	Gender    string    `json:"gender"     db:"gender"`
	Aadhaar   string    `json:"aadhaar"    db:"aadhaar"`
	Password  string    `json:"-"          db:"password"`
	Role      Role      `json:"role"       db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RiskLocation struct {
	ID        string    `json:"id"         db:"id"`
	Latitude  float64   `json:"latitude"   db:"latitude"`
	Longitude float64   `json:"longitude"  db:"longitude"`
	RiskLevel string    `json:"risk_level" db:"risk_level"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Report struct {
	ID          string    `json:"id"          db:"id"`
	ImageURL    string    `json:"image_url"   db:"image_url"`
	Latitude    float64   `json:"latitude"    db:"latitude"`
	Longitude   float64   `json:"longitude"   db:"longitude"`
	Description string    `json:"description" db:"description"`
	ReportedBy  string    `json:"reported_by" db:"reported_by"`
	Status      string    `json:"status"      db:"status"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
}

type SOS struct {
	ID        string    `json:"id"         db:"id"`
	UserID    string    `json:"user_id"    db:"user_id"`
	Latitude  float64   `json:"latitude"   db:"latitude"`
	Longitude float64   `json:"longitude"  db:"longitude"`
	Active    bool      `json:"active"     db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
