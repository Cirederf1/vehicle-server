package storage

import (
	"github.com/Cirederf1/vehicle-server/storage/vehiclestore"
)

type Store interface {
	Vehicle() vehiclestore.Store
}
