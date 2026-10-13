//go:build linux

package main

import (
	"errors"
	"syscall"
)

// kernelLacksIPv6 reports whether the kernel was built without IPv6, so
// that opening an AF_INET6 socket fails with EAFNOSUPPORT. Some sandboxes
// run such kernels, Claude Code's remote environments among them. Any other
// outcome, including another error, leaves Spanner Omni to run as it is.
func kernelLacksIPv6() bool {
	fd, err := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
	if err == nil {
		syscall.Close(fd)
		return false
	}
	return errors.Is(err, syscall.EAFNOSUPPORT)
}
