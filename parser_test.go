package strongdb

import (
	"path/filepath"
	"testing"
)

func TestParseSchema(t *testing.T) {
	s, err := ParseSchmea(filepath.Join("testdata", "equinox", "schema"))
	if err != nil {
		t.Error(err)
	}

	q, err := ParseQueries(filepath.Join("testdata", "equinox", "queries"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", q)

	if false {
		source := generate(s)
		t.Logf(source)
	}
}
