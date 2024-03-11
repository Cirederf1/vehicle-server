//go:build !integration

package vehicle_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cicd-lectures/vehicle-server/pkg/httputil"
	"github.com/cicd-lectures/vehicle-server/pkg/testutil"
	"github.com/cicd-lectures/vehicle-server/storage"
	"github.com/cicd-lectures/vehicle-server/vehicle"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateHandlerValidation(t *testing.T) {
	for _, testCase := range []struct {
		desc        string
		vehicle     vehicle.Vehicle
		wantStatus  int
		wantPayload httputil.APIError
	}{
		{
			desc:       "missing all fields",
			vehicle:    vehicle.Vehicle{},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"missing short code",
				},
			},
		},
		{
			desc: "short code too long",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aaabbbcccddd",
				Longitude:    44.3,
				Latitude:     23.4,
				BatteryLevel: 92,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"short code too long",
				},
			},
		},
		{
			desc: "negative battery level",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aabb",
				Longitude:    44.3,
				Latitude:     23.4,
				BatteryLevel: -92,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"battery level must be > 0 and <= 100",
				},
			},
		},
		{
			desc: "battery level above 100",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aabb",
				Longitude:    44.3,
				Latitude:     23.4,
				BatteryLevel: 130,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"battery level must be > 0 and <= 100",
				},
			},
		},
		{
			desc: "longitude below -90",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aabb",
				Longitude:    -93.3,
				Latitude:     23.4,
				BatteryLevel: 34,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"longitude must be >= -90 and <= 90",
				},
			},
		},
		{
			desc: "longitude above 90",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aabb",
				Longitude:    94.3,
				Latitude:     23.4,
				BatteryLevel: 34,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"longitude must be >= -90 and <= 90",
				},
			},
		},
		{
			desc: "latitude below -90",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aabb",
				Longitude:    33.3,
				Latitude:     -93.4,
				BatteryLevel: 34,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"latitude must be >= -90 and <= 90",
				},
			},
		},
		{
			desc: "latitude above 90",
			vehicle: vehicle.Vehicle{
				ShortCode:    "aabb",
				Longitude:    33.3,
				Latitude:     93.4,
				BatteryLevel: 34,
			},
			wantStatus: http.StatusBadRequest,
			wantPayload: httputil.APIError{
				Code:    httputil.ErrCodeInvalidRequestPayload,
				Message: "The request payload is invalid",
				Details: []any{
					"latitude must be >= -90 and <= 90",
				},
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			handler := vehicle.NewCreateHandler(
				storage.NewMemoryStore(),
				zap.NewNop(),
			)

			resp := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				"/vehicles",
				testutil.EncodeJSON(t, testCase.vehicle),
			)
			req.Header.Add("Content-Type", "application/json")

			handler.ServeHTTP(resp, req)

			assert.Equal(t, testCase.wantStatus, resp.Result().StatusCode)

			var gotPayload httputil.APIError
			httputil.DecodeJSON(resp.Result().Body, &gotPayload)
			assert.Equal(t, testCase.wantPayload, gotPayload)
		})
	}
}
