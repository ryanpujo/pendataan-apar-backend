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
	"github.com/ryanpujo/pendataan-apar-backend/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type KotaServiceMock struct {
	mock.Mock
}

func (ksm *KotaServiceMock) Create(ctx context.Context, payload *models.CreateKotaRequest) (uint, error) {
	args := ksm.Called(ctx, payload)
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return uint(args.Int(0)), args.Error(1)
}

func (ksm *KotaServiceMock) FindByID(ctx context.Context, id uint) (*models.Kota, error) {
	args := ksm.Called(ctx, id)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return args.Get(0).(*models.Kota), args.Error(1)
}

var (
	payload = &models.CreateKotaRequest{
		Name: "Jakarta",
	}
	kota = &models.Kota{
		ID:        1,
		Name:      "Jakarta",
		CreatedAt: now,
		UpdatedAt: now,
	}
)

func TestShouldCreateKota(t *testing.T) {
	jsonStr, _ := json.Marshal(payload)
	badPayload := models.CreateKotaRequest{
		Name: "",
	}
	badJson, _ := json.Marshal(badPayload)
	testTable := utils.TestTable{
		"success": {
			Payload: jsonStr,
			Arrange: func() {
				ksm.On("Create", mock.Anything, payload).Return(1, nil).Once()
			},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.CreateKotaResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusCreated, code)
				require.NotEmpty(t, res)
				require.Equal(t, uint(1), res.ID)
			},
		},
		"failed": {
			Payload: jsonStr,
			Arrange: func() {
				ksm.On("Create", mock.Anything, payload).Return(0, errors.New("failed to create kota")).Once()
			},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.CreateKotaResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, code)
				require.Empty(t, res.ID)
				require.Equal(t, "failed to create kota", res.Error)
			},
		},
		"validation error": {
			Payload: badJson,
			Arrange: func() {},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.CreateKotaResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, "validation error", res.Message)
				ksm.AssertNotCalled(t, "Create")
			},
		},
		"timeout cancel": {
			Payload: jsonStr,
			Arrange: func() {
				ksm.On("Create", mock.Anything, payload).After(time.Second*3).Return(1, nil).Once()
			},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.CreateKotaResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, context.DeadlineExceeded.Error(), res.Error)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.Arrange()

			req := httptest.NewRequest(http.MethodPost, "/v1/kota/create", bytes.NewReader(v.Payload))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var response controllers.CreateKotaResponse

			json.NewDecoder(res.Body).Decode(&response)

			v.AssertHTTP(t, res.Code, response)
		})
	}
}

func TestShouldFindKota(t *testing.T) {
	testTable := utils.TestTable{
		"success": {
			URIParam: 1,
			Arrange: func() {
				ksm.On("FindByID", mock.Anything, uint(1)).Return(kota, nil).Once()
			},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.FindKotaByIDResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusOK, code)
				require.Equal(t, *kota, res.Kota)
				ksm.AssertCalled(t, "FindByID", mock.Anything, uint(1))
			},
		},
		"failed": {
			URIParam: 1,
			Arrange: func() {
				ksm.On("FindByID", mock.Anything, uint(1)).Return((*models.Kota)(nil), errors.New("failed to find kota")).Once()
			},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.FindKotaByIDResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, code)
				require.Empty(t, res.Kota)
				ksm.AssertCalled(t, "FindByID", mock.Anything, uint(1))
			},
		},
		"validation error": {
			URIParam: 0,
			Arrange: func() {},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.FindKotaByIDResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, "validation error", res.Message)
			},
		},
		"timeout cancel": {
			URIParam: 1,
			Arrange: func() {
				ksm.On("FindByID", mock.Anything, uint(1)).After(time.Second*3).Return(kota, nil).Once()
			},
			AssertHTTP: func(t *testing.T, code int, response any) {
				res, ok := response.(controllers.FindKotaByIDResponse)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, context.DeadlineExceeded.Error(), res.Error)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.Arrange()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/kota/%d", v.URIParam), nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var response controllers.FindKotaByIDResponse

			json.NewDecoder(res.Body).Decode(&response)

			v.AssertHTTP(t, res.Code, response)
		})
	}
}
