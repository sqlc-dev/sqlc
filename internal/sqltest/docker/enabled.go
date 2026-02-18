package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sync/singleflight"
)

var flight singleflight.Group

func Installed() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found: %w", err)
	}
	return nil
}

// ensureContainer makes sure a Docker container with the given name is running.
// It handles three cases:
//  1. Container doesn't exist: create it with docker run
//  2. Container exists but stopped: start it with docker start
//  3. Container exists and running: do nothing
//
// It also handles the race condition where a parallel test process creates the
// container between our inspect and run calls.
func ensureContainer(name string, runArgs ...string) error {
	// Check if container exists and whether it's running
	output, err := exec.Command("docker", "container", "inspect",
		"-f", "{{.State.Running}}", name).CombinedOutput()

	if err == nil {
		// Container exists — check if it's running
		if strings.TrimSpace(string(output)) != "true" {
			// Container exists but is stopped, start it
			if startOut, startErr := exec.Command("docker", "start", name).CombinedOutput(); startErr != nil {
				return fmt.Errorf("docker start %s: %s: %w", name, string(startOut), startErr)
			}
		}
		return nil
	}

	// Container doesn't exist, create it
	args := append([]string{"run", "--name", name}, runArgs...)
	createOut, createErr := exec.Command("docker", args...).CombinedOutput()

	if createErr != nil {
		// Handle race: another process may have created the container between
		// our inspect and run calls
		conflictMsg := fmt.Sprintf(`Conflict. The container name "/%s" is already in use`, name)
		if strings.Contains(string(createOut), conflictMsg) {
			return nil
		}
		return fmt.Errorf("docker run %s: %s: %w", name, string(createOut), createErr)
	}

	return nil
}
