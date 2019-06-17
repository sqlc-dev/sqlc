package strongdb

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestParseSchema(t *testing.T) {
	s, err := ParseSchmea(filepath.Join("testdata", "ondeck", "schema"))
	if err != nil {
		t.Error(err)
	}

	q, err := ParseQueries(s, filepath.Join("testdata", "ondeck", "query"))
	if err != nil {
		t.Error(err)
	}

	source := generate(q, "ondeck")

	blob, err := ioutil.ReadFile(filepath.Join("testdata", "ondeck", "db.go"))
	if err != nil {
		log.Fatal(err)
	}

	if source != string(blob) {
		t.Errorf("output differs")
		fmt.Println(source)
	}
}
