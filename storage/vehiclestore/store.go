package vehiclestore

import "context"

type Point struct {
	Latitude  float64
	Longitude float64
}

type Vehicle struct {
	ID           int64
	ShortCode    string
	Position     Point
	BatteryLevel int64
}

type Store interface {
	// Creates a new vehicle.
	Create(context.Context, Vehicle) (Vehicle, error)

	// Finds the N closests vehicles from the current position.
	FindClosestFrom(context.Context, Point, int64) ([]Vehicle, error)

	// Delete a vehicle by its ID.
	// It returns true if the vehicle was deleted, false if the id did not exist.
	Delete(context.Context, int64) (bool, error)
}
