package dinosql

import (
	"io/ioutil"
)

func Exec(schemaDir, queryDir, pkg, out string, prepare bool) error {
	s, err := ParseSchmea(schemaDir)
	if err != nil {
		return err
	}

	q, err := ParseQueries(s, queryDir)
	if err != nil {
		return err
	}

	source := generate(q, pkg, prepare)
	return ioutil.WriteFile(out, []byte(source), 0644)
}
