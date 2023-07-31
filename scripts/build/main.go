package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	version := os.Getenv("VERSION")
	sha := os.Getenv("GITHUB_SHA")

	if version == "" {
		cmd := exec.Command("git", "show", "--no-patch", "--no-notes", "--pretty=%ci", sha)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(strings.TrimSpace(string(out)))
			log.Fatal(err)
		}
		var date string
		parts := strings.Split(string(out), " ")
		date = strings.Replace(parts[0]+parts[1], "-", "", -1)
		date = strings.Replace(date, ":", "", -1)
		version = fmt.Sprintf("v0.0.0-%s-%s", date, sha[:12])
	}

	fmt.Printf("::set-output name=version::%s\n", version)

	x := "-X github.com/sqlc-dev/sqlc/internal/cmd.version=" + version
	args := []string{
		"build",
		"-ldflags", x,
		"-o", "./sqlc",
		"./cmd/sqlc",
	}
	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(strings.TrimSpace(string(out)))
		log.Fatal(err)
	}
}
