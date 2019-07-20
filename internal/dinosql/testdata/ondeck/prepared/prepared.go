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
	if q.createCityStmt, err = db.PrepareContext(ctx, createCity); err != nil {
		return nil, err
	}
	if q.createVenueStmt, err = db.PrepareContext(ctx, createVenue); err != nil {
		return nil, err
	}
	if q.deleteVenueStmt, err = db.PrepareContext(ctx, deleteVenue); err != nil {
		return nil, err
	}
	if q.getCityStmt, err = db.PrepareContext(ctx, getCity); err != nil {
		return nil, err
	}
	if q.getVenueStmt, err = db.PrepareContext(ctx, getVenue); err != nil {
		return nil, err
	}
	if q.listCitiesStmt, err = db.PrepareContext(ctx, listCities); err != nil {
		return nil, err
	}
	if q.listVenuesStmt, err = db.PrepareContext(ctx, listVenues); err != nil {
		return nil, err
	}
	if q.updateCityNameStmt, err = db.PrepareContext(ctx, updateCityName); err != nil {
		return nil, err
	}
	if q.updateVenueNameStmt, err = db.PrepareContext(ctx, updateVenueName); err != nil {
		return nil, err
	}
	if q.venueCountByCityStmt, err = db.PrepareContext(ctx, venueCountByCity); err != nil {
		return nil, err
	}
	return &q, nil
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                   dbtx
	tx                   *sql.Tx
	createCityStmt       *sql.Stmt
	createVenueStmt      *sql.Stmt
	deleteVenueStmt      *sql.Stmt
	getCityStmt          *sql.Stmt
	getVenueStmt         *sql.Stmt
	listCitiesStmt       *sql.Stmt
	listVenuesStmt       *sql.Stmt
	updateCityNameStmt   *sql.Stmt
	updateVenueNameStmt  *sql.Stmt
	venueCountByCityStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                   tx,
		tx:                   tx,
		createCityStmt:       q.createCityStmt,
		createVenueStmt:      q.createVenueStmt,
		deleteVenueStmt:      q.deleteVenueStmt,
		getCityStmt:          q.getCityStmt,
		getVenueStmt:         q.getVenueStmt,
		listCitiesStmt:       q.listCitiesStmt,
		listVenuesStmt:       q.listVenuesStmt,
		updateCityNameStmt:   q.updateCityNameStmt,
		updateVenueNameStmt:  q.updateVenueNameStmt,
		venueCountByCityStmt: q.venueCountByCityStmt,
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

type CreateCityParams struct {
	Name string
	Slug string
}

func (q *Queries) CreateCity(ctx context.Context, arg CreateCityParams) (City, error) {
	row := q.queryRow(ctx, q.createCityStmt, createCity, arg.Name, arg.Slug)
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

type CreateVenueParams struct {
	Slug            string
	Name            string
	City            string
	SpotifyPlaylist string
	Status          Status
}

func (q *Queries) CreateVenue(ctx context.Context, arg CreateVenueParams) (int, error) {
	row := q.queryRow(ctx, q.createVenueStmt, createVenue, arg.Slug, arg.Name, arg.City, arg.SpotifyPlaylist, arg.Status)
	var id int
	err := row.Scan(&id)
	return id, err
}

const deleteVenue = `-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1 AND slug = $1
`

func (q *Queries) DeleteVenue(ctx context.Context, slug string) error {
	_, err := q.exec(ctx, q.deleteVenueStmt, deleteVenue, slug)
	return err
}

const getCity = `-- name: GetCity :one
SELECT slug, name
FROM city
WHERE slug = $1
`

func (q *Queries) GetCity(ctx context.Context, slug string) (City, error) {
	row := q.queryRow(ctx, q.getCityStmt, getCity, slug)
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const getVenue = `-- name: GetVenue :one
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, created_at
FROM venue
WHERE slug = $1 AND city = $2
`

type GetVenueParams struct {
	Slug string
	City string
}

func (q *Queries) GetVenue(ctx context.Context, arg GetVenueParams) (Venue, error) {
	row := q.queryRow(ctx, q.getVenueStmt, getVenue, arg.Slug, arg.City)
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
	rows, err := q.query(ctx, q.listCitiesStmt, listCities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []City
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
	rows, err := q.query(ctx, q.listVenuesStmt, listVenues, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Venue
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

type UpdateCityNameParams struct {
	Slug string
	Name string
}

func (q *Queries) UpdateCityName(ctx context.Context, arg UpdateCityNameParams) error {
	_, err := q.exec(ctx, q.updateCityNameStmt, updateCityName, arg.Slug, arg.Name)
	return err
}

const updateVenueName = `-- name: UpdateVenueName :one
UPDATE venue
SET name = $2
WHERE slug = $1
RETURNING id
`

type UpdateVenueNameParams struct {
	Slug string
	Name string
}

func (q *Queries) UpdateVenueName(ctx context.Context, arg UpdateVenueNameParams) (int, error) {
	row := q.queryRow(ctx, q.updateVenueNameStmt, updateVenueName, arg.Slug, arg.Name)
	var id int
	err := row.Scan(&id)
	return id, err
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
	rows, err := q.query(ctx, q.venueCountByCityStmt, venueCountByCity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []VenueCountByCityRow
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
