package vehicle

import "github.com/cicd-lectures/vehicle-server/storage/vehiclestore"

type Vehicle struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	ShortCode    string  `json:"shortcode"`
	BatteryLevel int64   `json:"battery"`
	ID           int64   `json:"id"`
}

func newVehicleFromModel(v vehiclestore.Vehicle) Vehicle {
	return Vehicle{
		ID:           v.ID,
		ShortCode:    v.ShortCode,
		Latitude:     v.Position.Latitude,
		Longitude:    v.Position.Longitude,
		BatteryLevel: v.BatteryLevel,
	}
}
