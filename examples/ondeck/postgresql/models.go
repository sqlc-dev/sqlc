// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package ondeck

import (
	"database/sql"
	"fmt"
	"time"
)

// Venues can be either open or closed
type Status string

const (
	StatusOpen   Status = "op!en"
	StatusClosed Status = "clo@sed"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type City struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Venues are places where muisc happens
type Venue struct {
	ID       int32    `json:"id"`
	Status   Status   `json:"status"`
	Statuses []Status `json:"statuses"`
	// This value appears in public URLs
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	City            string         `json:"city"`
	SpotifyPlaylist string         `json:"spotify_playlist"`
	SongkickID      sql.NullString `json:"songkick_id"`
	Tags            []string       `json:"tags"`
	CreatedAt       time.Time      `json:"created_at"`
}
