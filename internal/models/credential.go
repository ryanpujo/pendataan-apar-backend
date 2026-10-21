package models

import "time"

type Credential struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"-"` // Disembunyikan agar tidak bocor saat serialisasi JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RegisterCredentialRequest struct {
	Username    string `json:"username" binding:"required,min=4,max=50"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"` // Username/Email/Phone
	Password string `json:"password" validate:"required"`
}

type UpdateCredentialRequest struct {
	Username    *string `json:"username,omitempty"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Password    *string `json:"password,omitempty"`
}
