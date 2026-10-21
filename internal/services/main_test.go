package services_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/services"
)

var (
	crm               *CredentialRepoMock
	krm               *KotaRepoMock
	CredentialService services.Credential
	kotaService       services.Kota
	setupCTX          = func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), time.Second*1)
	}
	now = time.Now()
)

func TestMain(m *testing.M) {
	crm = new(CredentialRepoMock)
	krm = new(KotaRepoMock)

	CredentialService = services.NewCredentialService(crm)
	kotaService = services.NewKotaService(krm)

	os.Exit(m.Run())
}
