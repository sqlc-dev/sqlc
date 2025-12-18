package native

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Supported returns nil if native database installation is supported on this platform.
// Currently only Linux (Ubuntu/Debian) is supported.
func Supported() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("native database installation only supported on linux, got %s", runtime.GOOS)
	}
	// Check if apt-get is available (Debian/Ubuntu)
	if _, err := exec.LookPath("apt-get"); err != nil {
		return fmt.Errorf("apt-get not found: %w", err)
	}
	return nil
}
