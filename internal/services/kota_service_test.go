package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type KotaRepoMock struct {
	mock.Mock
}

func (krm *KotaRepoMock) Create(ctx context.Context, kota *models.CreateKotaRequest) (uint, error) {
	args := krm.Called(ctx, kota)
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return uint(args.Int(0)), args.Error(1)
}

func (krm *KotaRepoMock) FindByID(ctx context.Context, id uint) (*models.Kota, error) {
	args := krm.Called(ctx, id)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return args.Get(0).(*models.Kota), args.Error(1)
}

var (
	kotaPayload = &models.CreateKotaRequest{
		Name: "Jakarta Barat",
	}
	kota = &models.Kota{
		ID: 1,
		Name: "Jakarta Barat",
		CreatedAt: now,
		UpdatedAt: now,
	}
)

func TestShouldCreateKota(t *testing.T) {
	testTable := map[string]struct{
		setupCTX func() (context.Context, context.CancelFunc)
		arrange func(ctx context.Context)
		assert func(t *testing.T, actual uint, err error)
	}{
		"success": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				krm.On("Create", ctx, kotaPayload).Return(1, nil).Once()
			},
			assert: func(t *testing.T, actual uint, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, actual)
				require.Equal(t, uint(1), actual)
			},
		},
		"failed": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				krm.On("Create", ctx, kotaPayload).Return(0, errors.New("failed to create kota")).Once()
			},
			assert: func(t *testing.T, actual uint, err error) {
				require.Error(t, err)
				require.Zero(t, actual)
			},
		},
		"timeout cancel": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				krm.On("Create", ctx, kotaPayload).After(time.Second*3).Return(1, nil).Once()
			},
			assert: func(t *testing.T, actual uint, err error) {
				require.Error(t, err)
				require.Zero(t, actual)
				require.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := v.setupCTX()
			defer cancel()

			v.arrange(ctx)

			actual, err := kotaService.Create(ctx, kotaPayload)


			v.assert(t, actual, err)
		})
	}
}

func TestShouldFindKota(t *testing.T) {
	testTable := map[string]struct{
		setupCTX func() (context.Context, context.CancelFunc)
		arrange func(ctx context.Context)
		assert func(t *testing.T, actual *models.Kota, err error)
	}{
		"success": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				krm.On("FindByID", ctx, uint(1)).Return(kota, nil).Once()
			},
			assert: func(t *testing.T, actual *models.Kota, err error) {
				require.NoError(t, err)
				require.NotNil(t, actual)
				require.Equal(t, kota, actual)
			},
		},
		"failed": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				krm.On("FindByID", ctx, uint(1)).Return((*models.Kota)(nil), errors.New("failed to find kota")).Once()
			},
			assert: func(t *testing.T, actual *models.Kota, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.Equal(t, "failed to find kota", err.Error())
			},
		},
		"timeout cancel": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				krm.On("FindByID", ctx, uint(1)).After(time.Second*3).Return(kota, nil).Once()
			},
			assert: func(t *testing.T, actual *models.Kota, err error) {
				require.Error(t, err)
				require.Nil(t, actual)
				require.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}


	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := v.setupCTX()
			defer cancel()

			v.arrange(ctx)

			actual, err := kotaService.FindByID(ctx, 1)

			v.assert(t, actual, err)
		})
	}
}