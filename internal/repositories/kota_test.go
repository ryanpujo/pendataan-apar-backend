package repositories_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	kota = models.CreateKotaRequest{
		Name: "Jakarta",
	}
	created = models.Kota{
		ID:        1,
		Name:      "Jakarta",
		CreatedAt: now,
		UpdatedAt: now,
	}
)

func TestShouldCreateKota(t *testing.T) {
	testTable := map[string]struct {
		setupCTX func() (context.Context, context.CancelFunc)
		arrange  func()
		assert   func(t *testing.T, actual uint, err error)
	}{
		"success": {
			setupCTX: setupCTX,
			arrange: func() {
				mock.ExpectQuery("INSERT INTO kota").
					WithArgs(kota.Name, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			assert: func(t *testing.T, actual uint, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, actual)
				require.Equal(t, uint(1), actual)
			},
		},
		"failed": {
			setupCTX: setupCTX,
			arrange: func() {
				mock.ExpectQuery("INSERT INTO kota").
					WithArgs(kota.Name, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("failed to create new kota"))
			},
			assert: func(t *testing.T, actual uint, err error) {
				require.Error(t, err)
				require.Empty(t, actual)
				require.Equal(t, "failed to create new kota: failed to create new kota", err.Error())
			},
		},
		"timeout cancel": {
			setupCTX: func() (context.Context, context.CancelFunc) {
				ctx, cancel := setupCTX()
				cancel()
				return ctx, cancel
			},
			arrange: func() {

			},
			assert: func(t *testing.T, actual uint, err error) {
				require.Error(t, err)
				require.Empty(t, actual)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := v.setupCTX()
			defer cancel()

			v.arrange()

			actual, err := kotaRepoMock.Create(ctx, &kota)

			v.assert(t, actual, err)
		})
	}
}

func TestShouldFindKota(t *testing.T) {
	testTable := map[string]struct {
		setupCTX func() (context.Context, context.CancelFunc)
		arrange  func()
		assert   func(t *testing.T, actual *models.Kota, err error)
	}{
		"success": {
			setupCTX: setupCTX,
			arrange: func() {
				row := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).AddRow(
					1,
					kota.Name,
					now,
					now,
				)

				mock.ExpectQuery(`SELECT id, name, created_at, updated_at FROM kota WHERE id = \$1`).
					WithArgs(1).WillReturnRows(row)
			},
			assert: func(t *testing.T, actual *models.Kota, err error) {
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Equal(t, &created, actual)
			},
		},
		"failed": {
			setupCTX: setupCTX,
			arrange: func() {
				mock.ExpectQuery(`SELECT id, name, created_at, updated_at FROM kota WHERE id = \$1`).
					WithArgs(1).WillReturnError(errors.New("throw error"))
			},
			assert: func(t *testing.T, actual *models.Kota, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed to retrieve the kota: throw error", err.Error())
			},
		},
		"timeout cancel": {
			setupCTX: setupCTX,
			arrange: func() {
				row := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).AddRow(
					1,
					kota.Name,
					now,
					now,
				)

				mock.ExpectQuery(`SELECT id, name, created_at, updated_at FROM kota WHERE id = \$1`).
					WithArgs(1).WillDelayFor(time.Second * 3).WillReturnRows(row)
			},
			assert: func(t *testing.T, actual *models.Kota, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.ErrorIs(t, err, sqlmock.ErrCancelled)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			ctx, cancel := v.setupCTX()
			defer cancel()

			actual, err := kotaRepoMock.FindByID(ctx, 1)

			v.assert(t, actual, err)
		})
	}
}
