package repositories_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/pendataan-apar-backend/internal/repositories"
)

var (
	db                 *sql.DB
	mock               sqlmock.Sqlmock
	credentialRepoMock *repositories.CredentialRepo
	kotaRepoMock       *repositories.KotaRepo
	now                = time.Now().UTC()
	setupCTX           = func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), time.Second*1)
	}
)

func TestMain(m *testing.M) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	credentialRepoMock = repositories.NewCredentialRepo(db)
	kotaRepoMock = repositories.NewKotaRepo(db)

	os.Exit(m.Run())
}
