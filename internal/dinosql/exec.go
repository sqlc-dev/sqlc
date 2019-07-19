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

	c, err := ParseCatalog(settings.SchemaDir, settings)
	if err != nil {
		return err
	}

	q, err := ParseQueries(c, settings)
	if err != nil {
		return err
	}

	source, err := generate(q, settings)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(settings.Out, []byte(source), 0644)
}
