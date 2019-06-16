package ondeck

import (
	"context"
	"database/sql"
)

type City struct {
	Slug string
	Name string
}

type Venue struct {
	Slug            string
	Name            string
	City            sql.NullString
	SpotifyPlaylist string
	SongkickID      sql.NullString
}

// The shared methods on the sql.DB / sql.TX types
type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type DataLayer interface {
	GetCity(context.Context, string) (City, error)
	ListCities(context.Context) ([]City, error)
	ListVenues(context.Context, string) ([]Venue, error)
}

type Queries struct {
	db dbtx
}

const getCityQuery = `
SELECT slug, name
FROM city
WHERE slug = $1
`

func (q *Queries) GetCity(ctx context.Context, slug string) (City, error) {
	c := City{}
	err := q.db.QueryRowContext(ctx, getCityQuery, slug).Scan(&c.Slug, &c.Name)
	return c, err
}

const listCitiesQuery = `
SELECT slug, name
FROM city
ORDER BY name
`

func (q *Queries) ListCities(ctx context.Context) ([]City, error) {
	rows, err := q.db.QueryContext(ctx, listCitiesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []City{}
	for rows.Next() {
		c := City{}
		if err := rows.Scan(&c.Slug, &c.Name); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVenuesQuery = `
SELECT slug, name, city, spotify_playlist, songkick_id
FROM venue
WHERE city = $1
ORDER BY name
`

func (q *Queries) ListVenues(ctx context.Context, city string) ([]Venue, error) {
	rows, err := q.db.QueryContext(ctx, listVenuesQuery, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Venue{}
	for rows.Next() {
		v := Venue{}
		if err := rows.Scan(&v.Slug, &v.Name, &v.City, &v.SpotifyPlaylist, &v.SongkickID); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func New(d dbtx) *Queries {
	return &Queries{d}
}

type PreparedQueries struct {
	getCity    *sql.Stmt
	listCities *sql.Stmt
	listVenues *sql.Stmt

	tx *sql.Tx
}

func (q *PreparedQueries) WithTx(tx *sql.Tx) *PreparedQueries {
	return &PreparedQueries{
		getCity:    q.getCity,
		listCities: q.listCities,
		listVenues: q.listVenues,
		tx:         tx,
	}
}

func (q *PreparedQueries) GetCity(ctx context.Context, slug string) (City, error) {
	stmt := q.getCity
	if q.tx != nil {
		stmt = q.tx.StmtContext(ctx, stmt)
	}
	c := City{}
	err := stmt.QueryRowContext(ctx, slug).Scan(&c.Slug, &c.Name)
	return c, err
}

func (q *PreparedQueries) ListCities(ctx context.Context) ([]City, error) {
	stmt := q.listCities
	if q.tx != nil {
		stmt = q.tx.StmtContext(ctx, stmt)
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []City{}
	for rows.Next() {
		c := City{}
		if err := rows.Scan(&c.Slug, &c.Name); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *PreparedQueries) ListVenues(ctx context.Context, city string) ([]Venue, error) {
	stmt := q.listCities
	if q.tx != nil {
		stmt = q.tx.StmtContext(ctx, stmt)
	}
	rows, err := stmt.QueryContext(ctx, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Venue{}
	for rows.Next() {
		v := Venue{}
		if err := rows.Scan(&v.Slug, &v.Name, &v.City, &v.SpotifyPlaylist, &v.SongkickID); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func Prepare(ctx context.Context, db dbtx) (*PreparedQueries, error) {
	pq := PreparedQueries{}
	var err error
	if pq.getCity, err = db.PrepareContext(ctx, getCityQuery); err != nil {
		return nil, err
	}
	if pq.listCities, err = db.PrepareContext(ctx, listCitiesQuery); err != nil {
		return nil, err
	}
	if pq.listVenues, err = db.PrepareContext(ctx, listVenuesQuery); err != nil {
		return nil, err
	}
	return &pq, nil
}
