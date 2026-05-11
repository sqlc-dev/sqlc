package sqlc

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestGoModHasNoReplaceDirectives guards against regressions of
// https://github.com/sqlc-dev/sqlc/issues/4397.
//
// When go.mod contains a replace directive, the Go toolchain refuses to run
// `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` (and the equivalent
// `go run ...@latest`):
//
//	go: github.com/sqlc-dev/sqlc/cmd/sqlc@latest (in github.com/sqlc-dev/sqlc@v...):
//	    The go.mod file for the module providing named packages contains one or
//	    more replace directives. It must not contain directives that would cause
//	    it to be interpreted differently than if it were the main module.
//
// https://docs.sqlc.dev/en/latest/overview/install.html tells users to run
// exactly that command, so any replace directive slipping into go.mod breaks
// the advertised installation path for the next release.
func TestGoModHasNoReplaceDirectives(t *testing.T) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	var (
		inBlock   bool
		offenders []string
	)
	for i, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}

		if inBlock {
			if line == ")" {
				inBlock = false
				continue
			}
			if line != "" {
				offenders = append(offenders, fmt.Sprintf("  go.mod:%d: %s", i+1, raw))
			}
			continue
		}

		switch {
		case line == "replace (":
			inBlock = true
		case strings.HasPrefix(line, "replace "):
			offenders = append(offenders, fmt.Sprintf("  go.mod:%d: %s", i+1, raw))
		}
	}

	if len(offenders) > 0 {
		t.Fatalf("go.mod must not contain replace directives; "+
			"they break `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`.\n"+
			"See https://github.com/sqlc-dev/sqlc/issues/4397\n%s",
			strings.Join(offenders, "\n"))
	}
}
