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
	return nil
}
