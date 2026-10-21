package services

import (
	"context"

	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/ryanpujo/pendataan-apar-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type Credential interface {
	Register(ctx context.Context, payload *models.RegisterCredentialRequest) (uint, error)
	FindByID(ctx context.Context, id uint) (*models.Credential, error)
}

type CredentialService struct {
	CredRepo repositories.Credential
}

func NewCredentialService(credRepo repositories.Credential) *CredentialService {
	return &CredentialService{
		CredRepo: credRepo,
	}
}

var HashPassword = func(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (cs *CredentialService) Register(ctx context.Context, payload *models.RegisterCredentialRequest) (uint, error) {
	var err error
	payload.Password, err = HashPassword(payload.Password)
	if err != nil {
		return 0, err
	}
	return cs.CredRepo.Register(ctx, payload)
}

func (cs *CredentialService) FindByID(ctx context.Context, id uint) (*models.Credential, error) {
	return cs.CredRepo.FindByID(ctx, id)
}