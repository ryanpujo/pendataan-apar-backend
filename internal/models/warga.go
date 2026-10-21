package models

import "time"

type Warga struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	RT           string    `json:"rt"`
	RW           string    `json:"rw"`
	KelurahanID  int64     `json:"kelurahan_id"`
	CredentialID int64     `json:"credential_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ========================
// WARGA PAYLOADS
// ========================

// CreateWargaRequest untuk pembuatan profil Warga
type CreateWargaRequest struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	RT           string `json:"rt" validate:"required"`
	RW           string `json:"rw" validate:"required"`
	KelurahanID  int64  `json:"kelurahan_id" validate:"required"`
	CredentialID int64  `json:"credential_id" validate:"required"`
}

// CreateWargaWithCredentialsRequest jika registrasi warga mencakup pembuatan akun langsung
type CreateWargaWithCredentialsRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	RT          string `json:"rt" validate:"required"`
	RW          string `json:"rw" validate:"required"`
	KelurahanID int64  `json:"kelurahan_id" validate:"required"`
}

type UpdateWargaRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	RT          *string `json:"rt,omitempty"`
	RW          *string `json:"rw,omitempty"`
	KelurahanID *int64  `json:"kelurahan_id,omitempty"`
}