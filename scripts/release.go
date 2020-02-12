package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	flag.Parse()

	sha := os.Getenv("GITHUB_SHA")
	cmd := exec.Command("git", "show", "--no-patch", "--no-notes", "--pretty=%ci", sha)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	var date string
	parts := strings.Split(string(out), " ")
	date = strings.Replace(parts[0]+parts[1], "-", "", -1)
	date = strings.Replace(date, ":", "", -1)
	version := fmt.Sprintf("v0.0.0-%s-%s", date, sha[:12])

	x := "-X github.com/kyleconroy/sqlc/internal/cmd.version=" + version
	log.Printf("Releasing %s on channel %s", flag.Arg(1), flag.Arg(0))
	cmd = exec.Command("./equinox", "release",
		"--channel", flag.Arg(0),
		"--version", version,
		"--platforms", flag.Arg(1),
		"--app", "app_i4iCp1SuYfZ",
		"--token", os.Getenv("EQUINOX_API_TOKEN"),
		"--",
		"-ldflags", x, "./cmd/sqlc",
	)
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
}
