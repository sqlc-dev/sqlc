package bundler

import (
	"runtime"
	"time"

	"github.com/kyleconroy/sqlc/internal/info"
)

func projectMetadata() ([][2]string, error) {
	now := time.Now().UTC()
	return [][2]string{
		{"sqlc_version", info.Version},
		{"go_version", runtime.Version()},
		{"goos", runtime.GOOS},
		{"goarch", runtime.GOARCH},
		{"created", now.Format(time.RFC3339)},
	}, nil
}
