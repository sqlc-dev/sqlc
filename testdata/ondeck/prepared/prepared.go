package prepared

import (
	"context"
	"database/sql"
	"time"
)

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

type City struct {
	Slug string
	Name string
}

type Venue struct {
	ID              int
	Status          Status
	Slug            string
	Name            string
	City            string
	SpotifyPlaylist string
	SongkickID      sql.NullString
	CreatedAt       time.Time
}

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db dbtx) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createCity, err = db.PrepareContext(ctx, createCity); err != nil {
		return nil, err
	}
	if q.createVenue, err = db.PrepareContext(ctx, createVenue); err != nil {
		return nil, err
	}
	if q.deleteVenue, err = db.PrepareContext(ctx, deleteVenue); err != nil {
		return nil, err
	}
	if q.getCity, err = db.PrepareContext(ctx, getCity); err != nil {
		return nil, err
	}
	if q.getVenue, err = db.PrepareContext(ctx, getVenue); err != nil {
		return nil, err
	}
	if q.listCities, err = db.PrepareContext(ctx, listCities); err != nil {
		return nil, err
	}
	if q.listVenues, err = db.PrepareContext(ctx, listVenues); err != nil {
		return nil, err
	}
	if q.updateCityName, err = db.PrepareContext(ctx, updateCityName); err != nil {
		return nil, err
	}
	if q.updateVenueName, err = db.PrepareContext(ctx, updateVenueName); err != nil {
		return nil, err
	}
	if q.venueCountByCity, err = db.PrepareContext(ctx, venueCountByCity); err != nil {
		return nil, err
	}
	return &q, nil
}

type Queries struct {
	db               dbtx
	tx               *sql.Tx
	createCity       *sql.Stmt
	createVenue      *sql.Stmt
	deleteVenue      *sql.Stmt
	getCity          *sql.Stmt
	getVenue         *sql.Stmt
	listCities       *sql.Stmt
	listVenues       *sql.Stmt
	updateCityName   *sql.Stmt
	updateVenueName  *sql.Stmt
	venueCountByCity *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:               tx,
		tx:               tx,
		createCity:       q.createCity,
		createVenue:      q.createVenue,
		deleteVenue:      q.deleteVenue,
		getCity:          q.getCity,
		getVenue:         q.getVenue,
		listCities:       q.listCities,
		listVenues:       q.listVenues,
		updateCityName:   q.updateCityName,
		updateVenueName:  q.updateVenueName,
		venueCountByCity: q.venueCountByCity,
	}
}

const createCity = `-- name: CreateCity :one
INSERT INTO city (
    name,
    slug
) VALUES (
    $1,
    $2
) RETURNING slug, name
`

func (q *Queries) CreateCity(ctx context.Context, name string, slug string) (City, error) {
	var row *sql.Row
	switch {
	case q.createCity != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.createCity).QueryRowContext(ctx, name, slug)
	case q.createCity != nil:
		row = q.createCity.QueryRowContext(ctx, name, slug)
	default:
		row = q.db.QueryRowContext(ctx, createCity, name, slug)
	}
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const createVenue = `-- name: CreateVenue :one
INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status
) VALUES (
    $1,
    $2,
    $3,
    NOW(),
    $4,
    $5
) RETURNING id
`

func (q *Queries) CreateVenue(ctx context.Context, slug string, name string, city string, spotifyPlaylist string, status Status) (int, error) {
	var row *sql.Row
	switch {
	case q.createVenue != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.createVenue).QueryRowContext(ctx, slug, name, city, spotifyPlaylist, status)
	case q.createVenue != nil:
		row = q.createVenue.QueryRowContext(ctx, slug, name, city, spotifyPlaylist, status)
	default:
		row = q.db.QueryRowContext(ctx, createVenue, slug, name, city, spotifyPlaylist, status)
	}
	var i int
	err := row.Scan(&i)
	return i, err
}

const deleteVenue = `-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1 AND slug = $1
`

func (q *Queries) DeleteVenue(ctx context.Context, slug string) error {
	var err error
	switch {
	case q.deleteVenue != nil && q.tx != nil:
		_, err = q.tx.StmtContext(ctx, q.deleteVenue).ExecContext(ctx, slug)
	case q.deleteVenue != nil:
		_, err = q.deleteVenue.ExecContext(ctx, slug)
	default:
		_, err = q.db.ExecContext(ctx, deleteVenue, slug)
	}
	return err
}

const getCity = `-- name: GetCity :one
SELECT slug, name
FROM city
WHERE slug = $1
`

func (q *Queries) GetCity(ctx context.Context, slug string) (City, error) {
	var row *sql.Row
	switch {
	case q.getCity != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.getCity).QueryRowContext(ctx, slug)
	case q.getCity != nil:
		row = q.getCity.QueryRowContext(ctx, slug)
	default:
		row = q.db.QueryRowContext(ctx, getCity, slug)
	}
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const getVenue = `-- name: GetVenue :one
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, created_at
FROM venue
WHERE slug = $1 AND city = $2
`

func (q *Queries) GetVenue(ctx context.Context, slug string, city string) (Venue, error) {
	var row *sql.Row
	switch {
	case q.getVenue != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.getVenue).QueryRowContext(ctx, slug, city)
	case q.getVenue != nil:
		row = q.getVenue.QueryRowContext(ctx, slug, city)
	default:
		row = q.db.QueryRowContext(ctx, getVenue, slug, city)
	}
	var i Venue
	err := row.Scan(&i.ID, &i.Status, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID, &i.CreatedAt)
	return i, err
}

const listCities = `-- name: ListCities :many
SELECT slug, name
FROM city
ORDER BY name
`

func (q *Queries) ListCities(ctx context.Context) ([]City, error) {
	var rows *sql.Rows
	var err error
	switch {
	case q.listCities != nil && q.tx != nil:
		rows, err = q.tx.StmtContext(ctx, q.listCities).QueryContext(ctx)
	case q.listCities != nil:
		rows, err = q.listCities.QueryContext(ctx)
	default:
		rows, err = q.db.QueryContext(ctx, listCities)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []City{}
	for rows.Next() {
		var i City
		if err := rows.Scan(&i.Slug, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVenues = `-- name: ListVenues :many
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, created_at
FROM venue
WHERE city = $1
ORDER BY name
`

func (q *Queries) ListVenues(ctx context.Context, city string) ([]Venue, error) {
	var rows *sql.Rows
	var err error
	switch {
	case q.listVenues != nil && q.tx != nil:
		rows, err = q.tx.StmtContext(ctx, q.listVenues).QueryContext(ctx, city)
	case q.listVenues != nil:
		rows, err = q.listVenues.QueryContext(ctx, city)
	default:
		rows, err = q.db.QueryContext(ctx, listVenues, city)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Venue{}
	for rows.Next() {
		var i Venue
		if err := rows.Scan(&i.ID, &i.Status, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCityName = `-- name: UpdateCityName :exec
UPDATE city
SET name = $2
WHERE slug = $1
`

func (q *Queries) UpdateCityName(ctx context.Context, slug string, name string) error {
	var err error
	switch {
	case q.updateCityName != nil && q.tx != nil:
		_, err = q.tx.StmtContext(ctx, q.updateCityName).ExecContext(ctx, slug, name)
	case q.updateCityName != nil:
		_, err = q.updateCityName.ExecContext(ctx, slug, name)
	default:
		_, err = q.db.ExecContext(ctx, updateCityName, slug, name)
	}
	return err
}

const updateVenueName = `-- name: UpdateVenueName :one
UPDATE venue
SET name = $2
WHERE slug = $1
RETURNING id
`

func (q *Queries) UpdateVenueName(ctx context.Context, slug string, name string) (int, error) {
	var row *sql.Row
	switch {
	case q.updateVenueName != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.updateVenueName).QueryRowContext(ctx, slug, name)
	case q.updateVenueName != nil:
		row = q.updateVenueName.QueryRowContext(ctx, slug, name)
	default:
		row = q.db.QueryRowContext(ctx, updateVenueName, slug, name)
	}
	var i int
	err := row.Scan(&i)
	return i, err
}

const venueCountByCity = `-- name: VenueCountByCity :many
SELECT
    city,
    count(*)
FROM venue
GROUP BY 1
ORDER BY 1
`

type VenueCountByCityRow struct {
	City  string
	Count int
}

func (q *Queries) VenueCountByCity(ctx context.Context) ([]VenueCountByCityRow, error) {
	var rows *sql.Rows
	var err error
	switch {
	case q.venueCountByCity != nil && q.tx != nil:
		rows, err = q.tx.StmtContext(ctx, q.venueCountByCity).QueryContext(ctx)
	case q.venueCountByCity != nil:
		rows, err = q.venueCountByCity.QueryContext(ctx)
	default:
		rows, err = q.db.QueryContext(ctx, venueCountByCity)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []VenueCountByCityRow{}
	for rows.Next() {
		var i VenueCountByCityRow
		if err := rows.Scan(&i.City, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
