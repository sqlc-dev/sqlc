package strongdb

import (
	"path/filepath"
	"testing"
)

func TestParseSchema(t *testing.T) {
	s, err := ParseSchmea(filepath.Join("testdata", "ondeck", "schema"))
	if err != nil {
		t.Error(err)
	}

	q, err := ParseQueries(filepath.Join("testdata", "ondeck", "query"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", q)

	source := generate(s)
	t.Logf(source)
}
