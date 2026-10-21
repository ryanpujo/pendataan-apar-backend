package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
)

type Credential interface {
	Register(ctx context.Context, payload *models.RegisterCredentialRequest) (uint, error)
	FindByID(ctx context.Context, id uint) (*models.Credential, error)
}

type CredentialRepo struct {
	db *sql.DB
}

func NewCredentialRepo(db *sql.DB) *CredentialRepo {
	return &CredentialRepo{
		db: db,
	}
}

func (cr *CredentialRepo) Register(ctx context.Context, payload *models.RegisterCredentialRequest) (uint, error) {
	query := `
		INSERT INTO credentials (email, username, phone_number, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	now := time.Now().UTC()

	var id uint
	err := cr.db.QueryRowContext(ctx, query,
		payload.Email,
		payload.Username,
		payload.PhoneNumber,
		payload.Password,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to register credential: %w", err)
	}

	return id, nil
}

func (cr *CredentialRepo) FindByID(ctx context.Context, id uint) (*models.Credential, error) {
	query := `
		SELECT id, username, email, phone_number, password, created_at, updated_at
		FROM credentials WHERE id = $1
	`

	var credential models.Credential

	err := cr.db.QueryRowContext(ctx, query, id).Scan(
		&credential.ID,
		&credential.Username,
		&credential.Email,
		&credential.PhoneNumber,
		&credential.Password,
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no credential found by the id: %d: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve credential: %w", err)
	}

	return &credential, nil
}
