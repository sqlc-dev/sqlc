package docker

import (
	"fmt"
	"os/exec"

	"golang.org/x/sync/singleflight"
)

var flight singleflight.Group

func Installed() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found: %w", err)
	}
	// Verify the Docker daemon is actually running and accessible.
	// Without this check, tests will try Docker, fail on docker pull,
	// and t.Fatal instead of falling back to native database support.
	if out, err := exec.Command("docker", "info").CombinedOutput(); err != nil {
		return fmt.Errorf("docker daemon not available: %w\n%s", err, out)
	}
	return nil
}
