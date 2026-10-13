package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// v6to4Source is an LD_PRELOAD library that runs a program's IPv6 sockets
// over IPv4; see the comment at the top of the file.
//
//go:embed v6to4/v6to4.c
var v6to4Source []byte

// omniEnv is the environment Spanner Omni is started with. On a kernel with
// IPv6 it is this process's own. On one without, Omni's services abort at
// startup, as they open IPv6 sockets and have no flag to stop it ("Failed
// to create socket with ::socket ... Address family not supported by
// protocol" in their logs), so the v6to4 library is built into dir and
// preloaded into every Omni process, which inherit the environment from the
// launcher.
func omniEnv(dir string) ([]string, error) {
	env := os.Environ()
	if !kernelLacksIPv6() {
		return env, nil
	}
	log.Println("the kernel has no IPv6, which spanner omni needs: preloading v6to4 to run its IPv6 sockets over IPv4")
	lib, err := buildV6to4(dir)
	if err != nil {
		return nil, fmt.Errorf("building v6to4 for a kernel without IPv6: %w", err)
	}
	return append(env, "LD_PRELOAD="+lib), nil
}

// buildV6to4 compiles the v6to4 library into dir with the system's C
// compiler and returns its path.
func buildV6to4(dir string) (string, error) {
	if !commandExists("cc") {
		return "", fmt.Errorf("no C compiler (cc) is installed")
	}
	src := filepath.Join(dir, "v6to4.c")
	if err := os.WriteFile(src, v6to4Source, 0o644); err != nil {
		return "", err
	}
	lib := filepath.Join(dir, "libv6to4.so")
	// No fortified functions or stack protector, so the library needs no
	// glibc symbol newer than the 2.27 Omni ships with.
	err := run("cc", "-O2", "-shared", "-fPIC", "-fno-stack-protector",
		"-U_FORTIFY_SOURCE", "-D_FORTIFY_SOURCE=0", "-o", lib, src)
	if err != nil {
		return "", err
	}
	return lib, nil
}
