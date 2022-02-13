package bundler

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/kyleconroy/sqlc/internal/info"
)

type jsonTime struct {
	time.Time
}

func (t jsonTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))
	return []byte(stamp), nil
}

type setting struct {
	Key, Value string
}

type metadata struct {
	Version   string    `json:"version"`
	GoVersion string    `json:"go_version"`
	GOOS      string    `json:"goos"`
	GOARCH    string    `json:"goarch"`
	Created   jsonTime  `json:"created"`
	Settings  []setting `json:"settings"`
}

func buildMetadata() (*metadata, error) {
	return &metadata{
		Version:   info.Version,
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
		Created:   jsonTime{time.Now().UTC()},
	}, nil
}

func addMetadata(w *tar.Writer) error {
	md, err := buildMetadata()
	if err != nil {
		return err
	}
	blob, err := json.Marshal(md)
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name:       "metadata.json",
		Size:       int64(len(blob)),
		Mode:       0420,
		ModTime:    md.Created.Time,
		AccessTime: md.Created.Time,
		ChangeTime: md.Created.Time,
	}
	if err := w.WriteHeader(header); err != nil {
		return err
	}
	if _, err := w.Write(blob); err != nil {
		return err
	}
	return nil
}
