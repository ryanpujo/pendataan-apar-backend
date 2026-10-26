package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
)

type Kota interface {
	Create(ctx context.Context, kota *models.CreateKotaRequest) (uint, error)
	FindByID(ctx context.Context, id uint) (*models.Kota, error)
}

type KotaRepo struct {
	dB *sql.DB
}

func NewKotaRepo(db *sql.DB) *KotaRepo {
	return &KotaRepo{
		dB: db,
	}
}

func (kr *KotaRepo) Create(ctx context.Context, kota *models.CreateKotaRequest) (uint, error) {

	query := `
		INSERT INTO kota (name, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING id
	`

	now := time.Now().UTC()

	var id uint

	err := kr.dB.QueryRowContext(ctx, query,
		kota.Name,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create new kota: %w", err)
	}

	return id, nil
}

func (kr *KotaRepo) FindByID(ctx context.Context, id uint) (*models.Kota, error) {

	query := `
		SELECT id, name, created_at, updated_at FROM kota WHERE id = $1
	`

	var kota models.Kota

	err := kr.dB.QueryRowContext(ctx, query, id).Scan(
		&kota.ID,
		&kota.Name,
		&kota.CreatedAt,
		&kota.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kota by the id %d is not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve the kota: %w", err)
	}

	return &kota, nil
}