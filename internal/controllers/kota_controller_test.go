package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/controllers"
	"github.com/ryanpujo/pendataan-apar-backend/internal/models"
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

var (
	kota = &models.CreateKotaRequest{
		Name: "Jakarta",
	}
)

func TestShouldCreateKota(t *testing.T) {
	jsonStr, _ := json.Marshal(kota)
	badPayload := models.CreateKotaRequest{
		Name: "",
	}
	badJson, _ := json.Marshal(badPayload)
	testTable := map[string]struct {
		payload []byte
		arrange func()
		assert  func(t *testing.T, code int, res controllers.CreateKotaResponse)
	}{
		"success": {
			payload: jsonStr,
			arrange: func() {
				ksm.On("Create", mock.Anything, kota).Return(1, nil).Once()
			},
			assert: func(t *testing.T, code int, res controllers.CreateKotaResponse) {
				require.Equal(t, http.StatusCreated, code)
				require.NotEmpty(t, res)
				require.Equal(t, uint(1), res.ID)
			},
		},
		"failed": {
			payload: jsonStr,
			arrange: func() {
				ksm.On("Create", mock.Anything, kota).Return(0, errors.New("failed to create kota")).Once()
			},
			assert: func(t *testing.T, code int, res controllers.CreateKotaResponse) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Empty(t, res.ID)
				require.Equal(t, "failed to create kota", res.Error)
			},
		},
		"validation error": {
			payload: badJson,
			arrange: func() {},
			assert: func(t *testing.T, code int, res controllers.CreateKotaResponse) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, "validation error", res.Message)
				ksm.AssertNotCalled(t, "Create")
			},
		},
		"timeout cancel": {
			payload: jsonStr,
			arrange: func() {
				ksm.On("Create", mock.Anything, kota).After(time.Second*3).Return(1, nil).Once()
			},
			assert: func(t *testing.T, code int, res controllers.CreateKotaResponse) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, context.DeadlineExceeded.Error(), res.Error)
			},
		},
	}

	for k, v := range testTable {
		t.Run(k, func(t *testing.T) {
			v.arrange()

			req := httptest.NewRequest(http.MethodPost, "/v1/kota/create", bytes.NewReader(v.payload))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			var response controllers.CreateKotaResponse

			json.NewDecoder(res.Body).Decode(&response)

			v.assert(t, res.Code, response)
		})
	}
}
