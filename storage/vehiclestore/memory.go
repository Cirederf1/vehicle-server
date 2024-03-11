package vehiclestore

import (
	"context"
	"errors"
)

type MemoryStore struct {
	Data map[int64]Vehicle
	idx  int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{idx: 1, Data: make(map[int64]Vehicle)}
}

func (s *MemoryStore) Create(ctx context.Context, v Vehicle) (Vehicle, error) {
	v.ID = s.idx
	s.idx++

	s.Data[v.ID] = v

	return v, nil
}

func (s *MemoryStore) FindClosestFrom(ctx context.Context, location Point, limit int64) ([]Vehicle, error) {
	return nil, errors.New("not implemented")
}

func (s *MemoryStore) Delete(ctx context.Context, id int64) (bool, error) {
	return false, errors.New("not implemented")
}
