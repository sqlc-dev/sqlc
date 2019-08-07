package prepared

import (
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
	ID              int32
	Status          Status
	Slug            string
	Name            string
	City            string
	SpotifyPlaylist string
	SongkickID      sql.NullString
	Tags            []string
	CreatedAt       time.Time
}
