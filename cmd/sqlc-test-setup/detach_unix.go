//go:build unix

package main

import "syscall"

// detachedProcess puts a child in a session of its own, so it outlives
// this process and is not stopped by the signals this process gets.
func detachedProcess() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
