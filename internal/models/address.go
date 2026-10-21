package models

import "time"

// ========================
// ENTITY MODELS
// ========================

type Kota struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Kecamatan merepresentasikan tabel 'kecamatan'
type Kecamatan struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required"`
	KotaID    int64     `json:"kota_id" binding:"required,gt=0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Kelurahan merepresentasikan tabel 'kelurahan'
type Kelurahan struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	KodePos     string    `json:"kode_pos" binding:"required,len=5,numeric"`
	KecamatanID int64     `json:"kecamatan_id" binding:"required,gt=0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ========================
// KOTA PAYLOADS
// ========================

type CreateKotaRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateKotaRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// ========================
// KECAMATAN PAYLOADS
// ========================

type CreateKecamatanRequest struct {
	KotaID int64  `json:"kota_id" binding:"required,gt=0"`
	Name   string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateKecamatanRequest struct {
	KotaID *int64  `json:"kota_id,omitempty" binding:"omitempty,gt=0"`
	Name   *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
}

// ========================
// KELURAHAN PAYLOADS
// ========================

type CreateKelurahanRequest struct {
	KecamatanID int64  `json:"kecamatan_id" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	KodePos     string `json:"kode_pos" binding:"required,len=5,numeric"`
}

type UpdateKelurahanRequest struct {
	KecamatanID *int64  `json:"kecamatan_id,omitempty" binding:"omitempty,gt=0"`
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	KodePos     *string `json:"kode_pos,omitempty" binding:"omitempty,len=5,numeric"`
}