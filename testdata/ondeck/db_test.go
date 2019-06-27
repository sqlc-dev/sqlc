package ondeck

import (
	"context"
	"database/sql"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"

	"example.com/ondeck/prepared"

	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func id() string {
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func provision(t *testing.T, source string) (*sql.DB, func()) {
	t.Helper()

	db, err := sql.Open("postgres", source)
	if err != nil {
		t.Fatal(err)
	}

	schema := "dinotest_" + id()

	// For each test, pick a new schema name at random.
	// `foo` is used here only as an example
	if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}

	sdb, err := sql.Open("postgres", source+"&search_path="+schema)
	if err != nil {
		t.Fatal(err)
	}

	return sdb, func() {
		if _, err := db.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestQueries(t *testing.T) {
	t.Parallel()

	sdb, cleanup := provision(t, "postgres://localhost/dinotest?sslmode=disable")
	defer cleanup()

	files, err := ioutil.ReadDir("schema")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := ioutil.ReadFile(filepath.Join("schema", f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sdb.Exec(string(blob)); err != nil {
			t.Fatalf("%s: %s", f.Name(), err)
		}
	}

	q := New(sdb)

	ctx := context.Background()

	city, err := q.CreateCity(ctx, "san-francisco", "San Francisco")
	if err != nil {
		t.Fatal(err)
	}

	venueID, err := q.CreateVenue(ctx,
		"the-fillmore",
		"The Fillmore",
		city.Slug,
		"spotify:uri",
		StatusOpen)
	if err != nil {
		t.Fatal(err)
	}

	venue, err := q.GetVenue(ctx, "the-fillmore", city.Slug)
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
			{city.Slug, 1},
		}); diff != "" {
			t.Errorf("venue county mismatch:\n%s", diff)
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
		err := q.UpdateCityName(ctx, city.Slug, "SF")
		if err != nil {
			t.Error(err)
		}
	}

	{
		count, err := q.UpdateVenueName(ctx, venue.Slug, "Fillmore")
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(count, 1); diff != "" {
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

func TestPrepared(t *testing.T) {
	t.Parallel()

	sdb, cleanup := provision(t, "postgres://localhost/dinotest?sslmode=disable")
	defer cleanup()

	files, err := ioutil.ReadDir("schema")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		blob, err := ioutil.ReadFile(filepath.Join("schema", f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := sdb.Exec(string(blob)); err != nil {
			t.Fatalf("%s: %s", f.Name(), err)
		}
	}

	q, err := prepared.Prepare(context.Background(), sdb)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	city, err := q.CreateCity(ctx, "san-francisco", "San Francisco")
	if err != nil {
		t.Fatal(err)
	}

	venueID, err := q.CreateVenue(ctx,
		"the-fillmore",
		"The Fillmore",
		city.Slug,
		"spotify:uri",
		prepared.StatusOpen)
	if err != nil {
		t.Fatal(err)
	}

	venue, err := q.GetVenue(ctx, "the-fillmore", city.Slug)
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
		if diff := cmp.Diff(actual, []prepared.VenueCountByCityRow{
			{city.Slug, 1},
		}); diff != "" {
			t.Errorf("venue county mismatch:\n%s", diff)
		}
	}

	{
		actual, err := q.ListCities(ctx)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, []prepared.City{city}); diff != "" {
			t.Errorf("list city mismatch:\n%s", diff)
		}
	}

	{
		actual, err := q.ListVenues(ctx, city.Slug)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(actual, []prepared.Venue{venue}); diff != "" {
			t.Errorf("list venue mismatch:\n%s", diff)
		}
	}

	{
		err := q.UpdateCityName(ctx, city.Slug, "SF")
		if err != nil {
			t.Error(err)
		}
	}

	{
		count, err := q.UpdateVenueName(ctx, venue.Slug, "Fillmore")
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(count, 1); diff != "" {
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
