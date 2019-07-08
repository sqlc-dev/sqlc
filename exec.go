package dinosql

import (
	"io/ioutil"
)

func Exec(schemaDir, queryDir, out string, settings GenerateSettings) error {
	s, err := ParseSchmea(schemaDir)
	if err != nil {
		return err
	}

	q, err := ParseQueries(s, queryDir)
	if err != nil {
		return err
	}

	source := generate(q, settings)
	return ioutil.WriteFile(out, []byte(source), 0644)
}
