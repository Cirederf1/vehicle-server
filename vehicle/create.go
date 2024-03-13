package vehicle

import (
	"net/http"

	"github.com/Cirederf1/vehicle-server/pkg/httputil"
	"github.com/Cirederf1/vehicle-server/storage"
	"github.com/Cirederf1/vehicle-server/storage/vehiclestore"
	"go.uber.org/zap"
)

type CreateRequest struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	ShortCode    string  `json:"shortcode"`
	BatteryLevel int64   `json:"battery"`
}

func (f *CreateRequest) validate() []string {
	var validationIssues []string

	if f.ShortCode == "" {
		validationIssues = append(validationIssues, "missing short code")
	}

	if len(f.ShortCode) > 4 {
		validationIssues = append(validationIssues, "short code too long")
	}

	if f.Latitude < -90 || f.Latitude > 90 {
		validationIssues = append(validationIssues, "latitude must be >= -90 and <= 90")
	}

	if f.Longitude < -90 || f.Longitude > 90 {
		validationIssues = append(validationIssues, "longitude must be >= -90 and <= 90")
	}

	if f.BatteryLevel < 0 || f.BatteryLevel > 100 {
		validationIssues = append(validationIssues, "battery level must be > 0 and <= 100")
	}

	return validationIssues
}

type CreateResponse struct {
	Vehicle Vehicle `json:"vehicle"`
}

type CreateHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewCreateHandler(store storage.Store, logger *zap.Logger) *CreateHandler {
	return &CreateHandler{
		store:  store,
		logger: logger.With(zap.String("handler", "create_vehicle")),
	}
}

func (c *CreateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := httputil.DecodeRequestAsJSON(r, &req); err != nil {
		c.logger.Error(
			"Could not decode request body",
			zap.Error(err),
		)
		httputil.ServeError(rw, http.StatusBadRequest, err)
		return
	}

	if validationIssues := req.validate(); len(validationIssues) > 0 {
		httputil.ServeError(
			rw,
			http.StatusBadRequest,
			newValidationError(validationIssues),
		)
		return
	}

	newVehicle, err := c.store.Vehicle().Create(
		r.Context(),
		vehiclestore.Vehicle{
			ShortCode:    req.ShortCode,
			BatteryLevel: req.BatteryLevel,
			Position: vehiclestore.Point{
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
			},
		},
	)
	if err != nil {
		c.logger.Error(
			"Could not save the new vehicle",
			zap.Error(err),
		)
		httputil.ServeError(rw, http.StatusInternalServerError, err)
		return
	}

	httputil.ServeJSON(
		rw,
		http.StatusCreated,
		&CreateResponse{Vehicle: newVehicleFromModel(newVehicle)},
	)
}
