package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	payload = models.RegisterCredentialRequest{
		Username:    "ryanpujo",
		Email:       "ryanpujo@gmail.com",
		PhoneNumber: "085179812246",
		Password:    "okeoke",
	}
	credentialDummy = &models.Credential{
		ID:          1,
		Username:    "ryanpujo",
		Email:       "ryanpujo@gmail.com",
		PhoneNumber: "085179812246",
		Password:    "okeoke",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
)

func TestShouldCreateCredential(t *testing.T) {
	tableTest := map[string]struct {
		setupCTX func() (context.Context, context.CancelFunc)
		arrange  func()
		assert   func(t *testing.T, id uint, err error)
	}{
		"success": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*1)
			},
			arrange: func() {
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						payload.Email,
						payload.Username,
						payload.PhoneNumber,
						payload.Password,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Equal(t, uint(1), id)
				require.NoError(t, err)
			},
		},
		"duplicate email": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*1)
			},
			arrange: func() {
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						payload.Email,
						payload.Username,
						payload.PhoneNumber,
						payload.Password,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillReturnError(errors.New("email is already exists"))
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Zero(t, id)
				require.Error(t, err)
				require.Equal(t, "failed to register credential: email is already exists", err.Error())
			},
		},
		"timeout cancel": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
				return ctx, cancel
			},
			arrange: func() {
				mock.ExpectQuery("INSERT INTO credentials").
					WithArgs(
						payload.Email,
						payload.Username,
						payload.PhoneNumber,
						payload.Password,
						sqlmock.AnyArg(),
						sqlmock.AnyArg()).
					WillDelayFor(time.Second * 2).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Zero(t, id)
				require.Error(t, err)
			},
		},
	}

	for k, v := range tableTest {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			ctx, cancel := v.setupCTX()
			defer cancel()
			id, err := credentialRepoMock.Register(ctx, &payload)
			v.assert(t, id, err)
		})
	}
}

func TestShouldFindCredentialByID(t *testing.T) {
	testTable := map[string]struct {
		setupCTX func() (context.Context, context.CancelFunc)
		arrange  func(ctx context.Context)
		assert   func(t *testing.T, cred *models.Credential, err error)
	}{
		"success": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*1)
			},
			arrange: func(ctx context.Context) {
				rows := mock.NewRows([]string{"id", "username", "email", "phone_number", "password", "created_at", "updated_at"}).
					AddRow(1, payload.Username, payload.Email, payload.PhoneNumber, payload.Password, now, now)
				mock.ExpectQuery(`
					SELECT id, username. email, phone_number, password, created_at, updated_at
					FROM credentials WHERE id = \$1
				`).WithArgs(1).WillReturnRows(rows)
			},
			assert: func(t *testing.T, cred *models.Credential, err error) {
				require.NoError(t, err)
				require.Equal(t, credentialDummy, cred)
			},
		},
		"failed": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*1)
			},
			arrange: func(ctx context.Context) {
				mock.ExpectQuery(`
					SELECT id, username. email, phone_number, password, created_at, updated_at
					FROM credentials WHERE id = \$1
				`).WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			assert: func(t *testing.T, cred *models.Credential, err error) {
				require.Error(t, err)
				require.Nil(t, cred)
				require.ErrorIs(t, err, sql.ErrNoRows)
			},
		},
		"timeout cancel": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), time.Second*1)
			},
			arrange: func(ctx context.Context) {
				mock.ExpectQuery(`
					SELECT id, username. email, phone_number, password, created_at, updated_at
					FROM credentials WHERE id = \$1
				`).WithArgs(1).WillDelayFor(time.Second * 2).WillReturnError(sql.ErrNoRows)
			},
			assert: func(t *testing.T, cred *models.Credential, err error) {
				require.Error(t, err)
				require.Nil(t, cred)
				require.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := v.setupCTX()
			defer cancel()

			v.arrange(ctx)

			cred, err := credentialRepoMock.FindByID(ctx, 1)

			v.assert(t, cred, err)
		})
	}
}
