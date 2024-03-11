package vehiclestore

import (
	"context"
	"errors"

	pkgpgx "github.com/cicd-lectures/vehicle-server/pkg/pgx"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type PGXStore struct {
	conn pkgpgx.DB
}

func NewPGXStore(conn pkgpgx.DB) *PGXStore {
	return &PGXStore{conn: conn}
}

const createVehicleStatement = `
INSERT INTO vehicle_server.vehicles (shortcode, battery, position) VALUES ($1, $2, $3) RETURNING id;
`

func (p *PGXStore) Create(ctx context.Context, v Vehicle) (Vehicle, error) {
	encodedPos, err := ewkbhex.Encode(
		geom.NewPoint(geom.XY).
			MustSetCoords([]float64{v.Position.Longitude, v.Position.Latitude}).
			SetSRID(4326),
		ewkbhex.NDR,
	)
	if err != nil {
		return Vehicle{}, err
	}

	var id int64

	err = p.conn.QueryRow(
		ctx,
		createVehicleStatement,
		v.ShortCode,
		v.BatteryLevel,
		encodedPos,
	).Scan(&id)
	if err != nil {
		return Vehicle{}, err
	}

	return Vehicle{
		ID:           id,
		ShortCode:    v.ShortCode,
		BatteryLevel: v.BatteryLevel,
		Position:     v.Position,
	}, nil
}

const findClosestFromStatement = `
SELECT id, shortcode, battery, position
FROM vehicle_server.vehicles
ORDER BY position <-> ST_MakePoint($1, $2)::geography ASC
LIMIT $3;
`

var errInvalidCoordinates = errors.New("invalid coordinates")

func (p *PGXStore) FindClosestFrom(ctx context.Context, location Point, limit int64) ([]Vehicle, error) {
	var vehicles []Vehicle

	rows, err := p.conn.Query(ctx, findClosestFromStatement, location.Longitude, location.Latitude, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			v          Vehicle
			encodedPos string
		)

		if err := rows.Scan(
			&v.ID,
			&v.ShortCode,
			&v.BatteryLevel,
			&encodedPos,
		); err != nil {
			return nil, err
		}

		point, err := ewkbhex.Decode(encodedPos)
		if err != nil {
			return nil, err
		}
		coords := point.FlatCoords()
		if len(coords) != 2 {
			return nil, errInvalidCoordinates
		}

		v.Position.Longitude = coords[0]
		v.Position.Latitude = coords[1]

		vehicles = append(vehicles, v)
	}

	return vehicles, rows.Err()
}

const deleteByIDStatement = `
DELETE FROM vehicle_server.vehicles WHERE id = $1
`

func (p *PGXStore) Delete(ctx context.Context, id int64) (bool, error) {
	tag, err := p.conn.Exec(ctx, deleteByIDStatement, id)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() == 1, nil
}
