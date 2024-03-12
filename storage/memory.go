package storage

import "github.com/Cirederf1/vehicle-server/storage/vehiclestore"

type MemoryStore struct {
	VehicleStore *vehiclestore.MemoryStore
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{VehicleStore: vehiclestore.NewMemoryStore()}
}

func (m *MemoryStore) Vehicle() vehiclestore.Store {
	return m.VehicleStore
}
