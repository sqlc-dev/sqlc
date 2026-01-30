//go:build examples

package ondeck

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	_ "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

func runOnDeckQueries(t *testing.T, q *Queries) {
	ctx := context.Background()

	city, err := q.CreateCity(ctx, CreateCityParams{
		Slug: "san-francisco",
		Name: "San Francisco",
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	tags := "rock, punk"
	venueID, err := q.CreateVenue(ctx, CreateVenueParams{
		Slug:            "the-fillmore",
		Name:            "The Fillmore",
		City:            city.Slug,
		SpotifyPlaylist: "spotify:uri",
		Status:          "op!en",
		Tags:            &tags,
	}, query.WithIdempotent())
	if err != nil {
		t.Fatal(err)
	}

	venue, err := q.GetVenue(ctx, GetVenueParams{
		Slug: "the-fillmore",
		City: city.Slug,
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(venue.ID, venueID); diff != "" {
		t.Errorf("venue ID mismatch:\n%s", diff)
	}

	{
		actual, err := q.GetCity(ctx, city.Slug)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, city); diff != "" {
			t.Errorf("get city mismatch:\n%s", diff)
		}
	}

	{
		actual, err := q.VenueCountByCity(ctx)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, []VenueCountByCityRow{
			{city.Slug, uint64(1)},
		}); diff != "" {
			t.Errorf("venue count mismatch:\n%s", diff)
		}
	}

	{
		actual, err := q.ListCities(ctx)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, []City{city}); diff != "" {
			t.Errorf("list city mismatch:\n%s", diff)
		}
	}

	{
		actual, err := q.ListVenues(ctx, city.Slug)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, []Venue{venue}); diff != "" {
			t.Errorf("list venue mismatch:\n%s", diff)
		}
	}

	{
		err := q.UpdateCityName(ctx, UpdateCityNameParams{
			Slug: city.Slug,
			Name: "SF",
		}, query.WithIdempotent())
		if err != nil {
			t.Error(err)
		}
	}

	{
		id, err := q.UpdateVenueName(ctx, UpdateVenueNameParams{
			Slug: venue.Slug,
			Name: "Fillmore",
		}, query.WithIdempotent())
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(id, venue.ID); diff != "" {
			t.Errorf("update venue mismatch:\n%s", diff)
		}
	}

	{
		err := q.DeleteVenue(ctx, venue.Slug)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestQueries(t *testing.T) {
	ctx := context.Background()
	db := local.YDB(t, []string{"schema"})
	defer db.Close(ctx)

	runOnDeckQueries(t, New(db.Query()))
}
