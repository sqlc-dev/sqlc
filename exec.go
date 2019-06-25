package dinosql

import (
	"io/ioutil"
)

func Exec(schemaDir, queryDir, pkg, out string) error {
	s, err := ParseSchmea(schemaDir)
	if err != nil {
		return err
	}

	q, err := ParseQueries(s, queryDir)
	if err != nil {
		return err
	}

	source := generate(q, pkg)
	return ioutil.WriteFile(out, []byte(source), 0644)
}
