// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package ondeck

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateCity(ctx context.Context, arg CreateCityParams) error
	CreateVenue(ctx context.Context, arg CreateVenueParams) (sql.Result, error)
	DeleteVenue(ctx context.Context, arg DeleteVenueParams) error
	GetCity(ctx context.Context, slug string) (City, error)
	GetVenue(ctx context.Context, arg GetVenueParams) (Venue, error)
	ListCities(ctx context.Context) ([]City, error)
	ListVenues(ctx context.Context, city string) ([]Venue, error)
	UpdateCityName(ctx context.Context, arg UpdateCityNameParams) error
	UpdateVenueName(ctx context.Context, arg UpdateVenueNameParams) error
	VenueCountByCity(ctx context.Context) ([]VenueCountByCityRow, error)
}

var _ Querier = (*Queries)(nil)
