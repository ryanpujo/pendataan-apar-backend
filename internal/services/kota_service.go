package services

import (
	"context"

	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/ryanpujo/pendataan-apar-backend/internal/repositories"
)

type Kota interface {
	Create(ctx context.Context, payload *models.CreateKotaRequest) (uint, error)
	FindByID(ctx context.Context, id uint) (*models.Kota, error)
}

type KotaService struct {
	kotaRepo repositories.Kota
}

func NewKotaService(repo repositories.Kota) *KotaService {
	return &KotaService{
		kotaRepo: repo,
	}
}

func (ks *KotaService) Create(ctx context.Context, payload *models.CreateKotaRequest) (uint, error) {
	return ks.kotaRepo.Create(ctx, payload)
}

func (ks *KotaService) FindByID(ctx context.Context, id uint) (*models.Kota, error) {
	return ks.kotaRepo.FindByID(ctx, id)
}