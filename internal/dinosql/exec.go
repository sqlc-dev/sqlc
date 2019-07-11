package dinosql

import (
	"encoding/json"
	"io/ioutil"
)

func Exec(settingsPath string) error {
	blob, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		return err
	}

	var settings GenerateSettings
	if err := json.Unmarshal(blob, &settings); err != nil {
		return err
	}

	s, err := ParseSchmea(settings.SchemaDir, settings)
	if err != nil {
		return err
	}

	q, err := ParseQueries(s, settings.QueryDir)
	if err != nil {
		return err
	}

	source := generate(q, settings)
	return ioutil.WriteFile(settings.Out, []byte(source), 0644)
}
