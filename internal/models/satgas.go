package models

import "time"

type Satgas struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	KelurahanID  int64     `json:"kelurahan_id"`
	CredentialID int64     `json:"credential_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ========================
// SATGAS PAYLOADS
// ========================

// CreateSatgasRequest untuk registrasi profil Satgas baru
type CreateSatgasRequest struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	KelurahanID  int64  `json:"kelurahan_id" validate:"required"`
	CredentialID int64  `json:"credential_id" validate:"required"`
}

// CreateSatgasWithCredentialsRequest jika registrasi akun dan profil dilakukan dalam satu request
type CreateSatgasWithCredentialsRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	KelurahanID int64  `json:"kelurahan_id" validate:"required"`
}

type UpdateSatgasRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	KelurahanID *int64  `json:"kelurahan_id,omitempty"`
}