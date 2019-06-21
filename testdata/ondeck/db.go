package ondeck

import (
	"context"
	"database/sql"
	"time"
)

type City struct {
	Slug string
	Name string
}

type Venue struct {
	ID              int
	CreatedAt       time.Time
	Slug            string
	Name            string
	City            sql.NullString
	SpotifyPlaylist string
	SongkickID      sql.NullString
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
	return &q, nil
}

type Queries struct {
	db dbtx

	tx          *sql.Tx
	createCity  *sql.Stmt
	createVenue *sql.Stmt
	deleteVenue *sql.Stmt
	getCity     *sql.Stmt
	getVenue    *sql.Stmt
	listCities  *sql.Stmt
	listVenues  *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		tx:          tx,
		db:          tx,
		createCity:  q.createCity,
		createVenue: q.createVenue,
		deleteVenue: q.deleteVenue,
		getCity:     q.getCity,
		getVenue:    q.getVenue,
		listCities:  q.listCities,
		listVenues:  q.listVenues,
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
    name,
    slug,
    created_at,
    spotify_playlist,
    city
) VALUES (
    $1,
    $2,
    NOW(),
    $3,
    $4
) RETURNING id
`

func (q *Queries) CreateVenue(ctx context.Context, name string, slug string, spotify_playlist string, city string) (int, error) {
	var row *sql.Row
	switch {
	case q.createVenue != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.createVenue).QueryRowContext(ctx, name, slug, spotify_playlist, city)
	case q.createVenue != nil:
		row = q.createVenue.QueryRowContext(ctx, name, slug, spotify_playlist, city)
	default:
		row = q.db.QueryRowContext(ctx, createVenue, name, slug, spotify_playlist, city)
	}
	var i int
	err := row.Scan(&i)
	return i, err
}

const deleteVenue = `-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1
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
SELECT id, created_at, slug, name, city, spotify_playlist, songkick_id
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
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID)
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
SELECT id, created_at, slug, name, city, spotify_playlist, songkick_id
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
		if err := rows.Scan(&i.ID, &i.CreatedAt, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID); err != nil {
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
