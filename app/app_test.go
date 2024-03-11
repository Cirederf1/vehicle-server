//go:build integration

package app_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/cicd-lectures/vehicle-server/pkg/httputil"
	"github.com/cicd-lectures/vehicle-server/pkg/testutil"
	"github.com/cicd-lectures/vehicle-server/storage/vehiclestore"
	"github.com/cicd-lectures/vehicle-server/vehicle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_CreatesVehicles(t *testing.T) {
	t.Parallel()

	// Setup the testenvironment, and clean it up as soon as the test finishes.
	app, teardown := setupEnvironment(t)
	t.Cleanup(teardown)

	// Creates a new vehicle, and make assertions on the result.
	newVehicle := vehicle.CreateRequest{
		Latitude:     10.0,
		Longitude:    9.0,
		ShortCode:    "ebvf",
		BatteryLevel: 72,
	}

	resp, err := http.Post(
		"http://"+app.ListenAddress()+"/vehicles",
		"application/json",
		testutil.EncodeJSON(t, &newVehicle),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	// Decode the response body, and make assertions on the body.
	var (
		gotResponse  vehicle.CreateResponse
		wantResponse = vehicle.CreateResponse{
			Vehicle: vehicle.Vehicle{
				ID:           1,
				Latitude:     newVehicle.Latitude,
				Longitude:    newVehicle.Longitude,
				ShortCode:    newVehicle.ShortCode,
				BatteryLevel: newVehicle.BatteryLevel,
			},
		}
	)

	err = httputil.DecodeJSON(resp.Body, &gotResponse)
	require.NoError(t, err)
	assert.Equal(t, wantResponse, gotResponse)
}

var vehicleSeed = []vehiclestore.Vehicle{
	{Position: vehiclestore.Point{Latitude: 50.0, Longitude: 50.0}, ShortCode: "aaa", BatteryLevel: 40},
	{Position: vehiclestore.Point{Latitude: 51.0, Longitude: 51.0}, ShortCode: "bbb", BatteryLevel: 50},
	{Position: vehiclestore.Point{Latitude: 52.0, Longitude: 52.0}, ShortCode: "ccc", BatteryLevel: 60},
}

func TestApp_ListsClosestVehicles(t *testing.T) {
	t.Parallel()
	// Setup the testenvironment, and clean it up as soon as the test finishes.
	app, teardown := setupEnvironment(t)
	t.Cleanup(teardown)

	// Add some vehicles to the database.
	seedVehicles(
		t,
		app.Store().Vehicle(),
		vehicleSeed...,
	)

	// Make a request.
	resp, err := http.Get(
		"http://" + app.ListenAddress() + "/vehicles?latitude=49.0&longitude=49.0&limit=10",
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	// Decode the response body, and make assertions on the body.
	var (
		gotResponse  vehicle.ListResponse
		wantResponse = vehicle.ListResponse{
			Vehicles: []vehicle.Vehicle{
				{ID: 1, Latitude: 50.0, Longitude: 50.0, ShortCode: "aaa", BatteryLevel: 40},
				{ID: 2, Latitude: 51.0, Longitude: 51.0, ShortCode: "bbb", BatteryLevel: 50},
				{ID: 3, Latitude: 52.0, Longitude: 52.0, ShortCode: "ccc", BatteryLevel: 60},
			},
		}
	)

	err = httputil.DecodeJSON(resp.Body, &gotResponse)
	require.NoError(t, err)
	assert.Equal(t, wantResponse, gotResponse)
}

func TestApp_DeletesVehicleByID(t *testing.T) {
	// Setup the testenvironment, and clean it up as soon as the test finishes.
	app, teardown := setupEnvironment(t)
	t.Cleanup(teardown)

	// Add some vehicles to the database.
	seedVehicles(
		t,
		app.Store().Vehicle(),
		vehicleSeed...,
	)

	// Make a request for a vehicle that does not exists.
	// We should get a 404.
	resp, err := httpDelete(
		"http://" + app.ListenAddress() + "/vehicles/456",
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Make a request for a vehicle that does exist.
	// We should get a  204.
	resp, err = httpDelete(
		"http://" + app.ListenAddress() + "/vehicles/1",
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// List all vehicles, we should have only two.
	vehicles, err := app.Store().Vehicle().FindClosestFrom(
		context.Background(),
		vehiclestore.Point{
			Longitude: 10.0,
			Latitude:  10.0,
		},
		10,
	)
	require.NoError(t, err)
	assert.Equal(t,
		[]vehiclestore.Vehicle{
			{
				ID:        2,
				ShortCode: "bbb",
				Position: vehiclestore.Point{
					Latitude:  51,
					Longitude: 51,
				},
				BatteryLevel: 50,
			},
			{
				ID:        3,
				ShortCode: "ccc",
				Position: vehiclestore.Point{
					Latitude:  52,
					Longitude: 52,
				},
				BatteryLevel: 60,
			},
		},
		vehicles,
	)
}

func httpDelete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
