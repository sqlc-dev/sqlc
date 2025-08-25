package docker

import (
	"fmt"
	"os/exec"
)

func Installed() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found: %w", err)
	}
	return nil
}
