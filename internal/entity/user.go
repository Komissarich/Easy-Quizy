package entity

import (
	"errors"
)

type User struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  Role   `json:"role"`

	PassHash []byte
}

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
	// SuperAdminRole Role = "superadmin" // ??? do we need it ???
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
