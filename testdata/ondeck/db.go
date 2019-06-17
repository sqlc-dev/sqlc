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

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Queries struct {
	db dbtx

	tx         *sql.Tx
	listCities *sql.Stmt
	getCity    *sql.Stmt
	listVenues *sql.Stmt
}

const listCities = `
SELECT slug, name FROM city WHERE  ORDER BY name
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
		i := City{}
		if err := row.Scan(&i.Slug, &i.Name); err != nil {
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

const getCity = `
SELECT slug, name FROM city WHERE slug = $1
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
	i := City{}
	err := row.Scan(&i.Slug, &i.Name)
	return c, err
}

const listVenues = `
SELECT slug, name, city, spotify_playlist, songkick_id FROM venue WHERE city = $1 ORDER BY name
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
		i := Venue{}
		if err := row.Scan(&i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID); err != nil {
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
