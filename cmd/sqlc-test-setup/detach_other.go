//go:build !unix

package main

import "syscall"

func detachedProcess() *syscall.SysProcAttr {
	return nil
}
