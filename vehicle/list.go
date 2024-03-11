package vehicle

import (
	"net/http"
	"strconv"

	"github.com/cicd-lectures/vehicle-server/pkg/httputil"
	"github.com/cicd-lectures/vehicle-server/storage"
	"github.com/cicd-lectures/vehicle-server/storage/vehiclestore"
	"go.uber.org/zap"
)

type ListRequest struct {
	Latitude  float64
	Longitude float64
	Limit     int64
}

func newListRequestFromQueryParameters(r *http.Request) *ListRequest {
	var req ListRequest

	req.Latitude, _ = strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	req.Longitude, _ = strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	req.Limit, _ = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	return &req
}

type ListResponse struct {
	Vehicles []Vehicle `json:"vehicles"`
}

func newListResponse(vehicles []vehiclestore.Vehicle) *ListResponse {
	result := make([]Vehicle, len(vehicles))

	for i, v := range vehicles {
		result[i] = newVehicleFromModel(v)
	}

	return &ListResponse{Vehicles: result}
}

type ListHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewListHandler(store storage.Store, logger *zap.Logger) *ListHandler {
	return &ListHandler{
		store:  store,
		logger: logger.With(zap.String("handler", "list_vehicles")),
	}
}

func (l *ListHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	req := newListRequestFromQueryParameters(r)

	vehicles, err := l.store.Vehicle().FindClosestFrom(
		r.Context(),
		vehiclestore.Point{Latitude: req.Latitude, Longitude: req.Longitude},
		req.Limit,
	)
	if err != nil {
		l.logger.Error(
			"Could not list vehicles from store",
			zap.Error(err),
		)

		httputil.ServeError(rw, http.StatusInternalServerError, err)
		return
	}

	httputil.ServeJSON(rw, http.StatusOK, newListResponse(vehicles))
}
