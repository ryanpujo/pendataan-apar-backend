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

type CredentialRepoMock struct {
	mock.Mock
}

func (cm *CredentialRepoMock) Register(ctx context.Context, payload *models.RegisterCredentialRequest) (uint, error) {
	args := cm.Called(ctx, payload)
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return uint(args.Int(0)), args.Error(1)
}

func (cm *CredentialRepoMock) FindByID(ctx context.Context, id uint) (*models.Credential, error) {
	args := cm.Called(ctx, id)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return args.Get(0).(*models.Credential), args.Error(1)
}

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
	testTable := map[string]struct {
		setupCTX func() (context.Context, context.CancelFunc)
		arrange  func(ctx context.Context)
		assert   func(t *testing.T, id uint, err error)
	}{
		"success": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				crm.On("Register", ctx, &payload).Return(1, nil).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Equal(t, uint(1), id)
				require.NoError(t, err)
			},
		},
		"failed": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				crm.On("Register", ctx, &payload).Return(0, errors.New("failed to create credential")).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Zero(t, id)
				require.Error(t, err)
			},
		},
		"timeout cancel": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				crm.On("Register", ctx, &payload).After(time.Second*2).Return(1, nil).Once()
			},
			assert: func(t *testing.T, id uint, err error) {
				require.Zero(t, 0)
				require.Error(t, err)
				require.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			ctx, cancel := v.setupCTX()
			defer cancel()

			v.arrange(ctx)

			id, err := CredentialService.Register(ctx, &payload)

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
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				crm.On("FindByID", ctx, uint(1)).Return(credentialDummy, nil).Once()
			},
			assert: func(t *testing.T, cred *models.Credential, err error) {
				require.NoError(t, err)
				require.Equal(t, credentialDummy, cred)
			},
		},
		"failed": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				crm.On("FindByID", ctx, uint(1)).Return((*models.Credential)(nil), errors.New("failed to create credential")).Once()
			},
			assert: func(t *testing.T, cred *models.Credential, err error) {
				require.Error(t, err)
				require.Nil(t, cred)
			},
		},
		"timeout cancel": {
			setupCTX: setupCTX,
			arrange: func(ctx context.Context) {
				crm.On("FindByID", ctx, uint(1)).After(time.Second*2).Return(credentialDummy, nil).Once()
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

			actual, err := CredentialService.FindByID(ctx, 1)

			v.assert(t, actual, err)
		})
	}
}
