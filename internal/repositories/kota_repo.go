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
		VALUES ($1, $2, $3) return id
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