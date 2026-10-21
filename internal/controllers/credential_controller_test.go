package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/controllers"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ServiceMock struct {
	mock.Mock
}

func (sm *ServiceMock) Register(ctx context.Context, payload *models.RegisterCredentialRequest) (uint, error) {
	args := sm.Called(ctx, payload)

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	return uint(args.Int(0)), args.Error(1)
}

func (sm *ServiceMock) FindByID(ctx context.Context, id uint) (*models.Credential, error) {
	args := sm.Called(ctx, id)

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return args.Get(0).(*models.Credential), args.Error(1)
}

var (
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
	payload := models.RegisterCredentialRequest{
		Username:    "ryanpujo",
		Email:       "ryanpujo@gmail.com",
		PhoneNumber: "085179812246",
		Password:    "okeokeoke",
	}
	invalidPayload := models.RegisterCredentialRequest{
		Username:    "df",
		Email:       "ryanpujo@gmail.com",
		PhoneNumber: "085179812246",
		Password:    "okeokeoke",
	}
	jsonPayload, _ := json.Marshal(payload)
	invalidJson, _ := json.Marshal(invalidPayload)
	testTable := map[string]struct {
		payload []byte
		arrange func()
		assert  func(t *testing.T, statusCode int, res controllers.RegisterResponse)
	}{
		"success": {
			payload: jsonPayload,
			arrange: func() {
				sm.On("Register", mock.Anything, &payload).Return(1, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, res controllers.RegisterResponse) {
				require.Equal(t, http.StatusCreated, statusCode)
				require.NotEmpty(t, res)
				require.Equal(t, uint(1), res.ID)
			},
		},
		"failed": {
			payload: jsonPayload,
			arrange: func() {
				sm.On("Register", mock.Anything, &payload).Return(0, errors.New("failed to create credential")).Once()
			},
			assert: func(t *testing.T, statusCode int, res controllers.RegisterResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotEmpty(t, res)
				require.Equal(t, "failed to create credential", res.Message)
			},
		},
		"invalid json": {
			payload: invalidJson,
			arrange: func() {},
			assert: func(t *testing.T, statusCode int, res controllers.RegisterResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotEmpty(t, res)
				require.Equal(t, "validation error", res.Message)
			},
		},
		"timeout cancel": {
			payload: jsonPayload,
			arrange: func() {
				sm.On("Register", mock.Anything, &payload).After(time.Second*2).Return(1, nil).Once()
			},
			assert: func(t *testing.T, statusCode int, res controllers.RegisterResponse) {
				require.Equal(t, http.StatusBadRequest, statusCode)
				require.NotEmpty(t, res)
				require.Equal(t, "request timeout", res.Message)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/credentials/register", v1), bytes.NewReader(v.payload))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var jsonRes controllers.RegisterResponse

			json.NewDecoder(res.Body).Decode(&jsonRes)

			v.assert(t, res.Code, jsonRes)
		})
	}
}

func TestShouldFindCredentialByID(t *testing.T) {
	testTable := map[string]struct {
		id      uint
		arrange func()
		assert  func(t *testing.T, code int, res controllers.FindByIDResponse)
	}{
		"success": {
			id: 1,
			arrange: func() {
				sm.On("FindByID", mock.Anything, uint(1)).Return(credentialDummy, nil).Once()
			},
			assert: func(t *testing.T, code int, res controllers.FindByIDResponse) {
				credentialDummy = &models.Credential{
					ID:          1,
					Username:    "ryanpujo",
					Email:       "ryanpujo@gmail.com",
					PhoneNumber: "085179812246",
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				require.Equal(t, http.StatusOK, code)
				require.Equal(t, credentialDummy, res.Credential)
			},
		},
		"failed": {
			id: 1,
			arrange: func() {
				sm.On("FindByID", mock.Anything, uint(1)).Return((*models.Credential)(nil), errors.New("failed to find credential")).Once()
			},
			assert: func(t *testing.T, code int, res controllers.FindByIDResponse) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, "failed to find credential", res.Error)
			},
		},
		"timeout cancel": {
			id: 1,
			arrange: func() {
				sm.On("FindByID", mock.Anything, uint(1)).After(time.Second*2).Return(credentialDummy, nil).Once()
			},
			assert: func(t *testing.T, code int, res controllers.FindByIDResponse) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, context.DeadlineExceeded.Error(), res.Error)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/credentials/%d", v1, v.id), nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var resJSON controllers.FindByIDResponse

			json.NewDecoder(res.Body).Decode(&resJSON)

			v.assert(t, res.Code, resJSON)
		})
	}
}
