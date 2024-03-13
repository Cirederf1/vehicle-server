package vehicle

import (
	"net/http"
	"strconv"

	"github.com/Cirederf1/vehicle-server/storage"
	"go.uber.org/zap"
)

type DeleteHandler struct {
	store  storage.Store
	logger *zap.Logger
}

func NewDeleteHandler(store storage.Store, logger *zap.Logger) *DeleteHandler {
	return &DeleteHandler{
		store:  store,
		logger: logger.With(zap.String("handler", "delete_vehicles")),
	}
}

func (d *DeleteHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err!=nil{
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	check, err := d.store.Vehicle().Delete(r.Context(), id)

	if err != nil {
		http.Error(rw, "Test", http.StatusInternalServerError)
		return
	}
	if check {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
}
