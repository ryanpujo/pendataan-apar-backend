package repositories_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	kota = models.CreateKotaRequest{
		Name: "Jakarta",
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
