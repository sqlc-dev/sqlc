//go:build examples

package ondeck

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest"

	"github.com/google/go-cmp/cmp"
)

// TODO: Enum is not yet supported
const (
	VenuesStatusOpen   string = "open"
	VenuesStatusClosed string = "closed"
)

func join(vals ...string) sql.NullString {
	if len(vals) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		Valid:  true,
		String: strings.Join(vals, ","),
	}
}

func runOnDeckQueries(t *testing.T, q *Queries) {
	ctx := context.Background()

	err := q.CreateCity(ctx, CreateCityParams{
		Slug: "san-francisco",
		Name: "San Francisco",
	})
	if err != nil {
		t.Fatal(err)
	}

	city, err := q.GetCity(ctx, "san-francisco")
	if err != nil {
		t.Fatal(err)
	}

	venueResult, err := q.CreateVenue(ctx, CreateVenueParams{
		Slug:            "the-fillmore",
		Name:            "The Fillmore",
		City:            city.Slug,
		SpotifyPlaylist: "spotify:uri",
		Status:          VenuesStatusOpen,
		Statuses:        join(string(VenuesStatusOpen), string(VenuesStatusClosed)),
		Tags:            join("rock", "punk"),
	})
	if err != nil {
		t.Fatal(err)
	}
	venueID, err := venueResult.LastInsertId()
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
		actual, err := q.VenueCountByCity(ctx)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, []VenueCountByCityRow{
			{city.Slug, int64(1)},
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
		})
		if err != nil {
			t.Error(err)
		}
	}

	{
		expected := "Fillmore"
		err := q.UpdateVenueName(ctx, UpdateVenueNameParams{
			Slug: venue.Slug,
			Name: expected,
		})
		if err != nil {
			t.Error(err)
		}
		fresh, err := q.GetVenue(ctx, GetVenueParams{
			Slug: venue.Slug,
			City: city.Slug,
		})
		if diff := cmp.Diff(expected, fresh.Name); diff != "" {
			t.Errorf("update venue mismatch:\n%s", diff)
		}
	}

	{
		err := q.DeleteVenue(ctx, DeleteVenueParams{
			Slug:   venue.Slug,
			Slug_2: venue.Slug,
		})
		if err != nil {
			t.Error(err)
		}
	}
}

func TestPrepared(t *testing.T) {
	t.Parallel()

	sdb, cleanup := sqltest.SQLite(t, []string{"schema"})
	defer sdb.Close()
	defer cleanup()

	q, err := Prepare(context.Background(), sdb)
	if err != nil {
		t.Fatal(err)
	}

	runOnDeckQueries(t, q)
}

func TestQueries(t *testing.T) {
	t.Parallel()

	sdb, cleanup := sqltest.SQLite(t, []string{"schema"})
	defer sdb.Close()
	defer cleanup()

	runOnDeckQueries(t, New(sdb))
}
