package storage

import (
	"github.com/cicd-lectures/vehicle-server/storage/vehiclestore"
)

type Store interface {
	Vehicle() vehiclestore.Store
}
