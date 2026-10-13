//go:build !linux

package main

// kernelLacksIPv6 is false off Linux, where Spanner Omni is not published.
func kernelLacksIPv6() bool {
	return false
}
