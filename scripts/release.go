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
	draft := flag.Bool("draft", false, "create a draft release")
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

	xname := "./equinox"
	if _, err := os.Stat("./equinox"); os.IsNotExist(err) {
		xname = "equinox"
	}

	args := []string{"release",
		"--channel", flag.Arg(0),
		"--version", version,
	}

	if *draft {
		args = append(args, "--draft")
	}

	args = append(args, []string{
		"--platforms", flag.Arg(1),
		"--app", "app_i4iCp1SuYfZ",
		"--token", os.Getenv("EQUINOX_API_TOKEN"),
		"--",
		"-ldflags", x, "./cmd/sqlc",
	}...)

	cmd = exec.Command(xname, args...)
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	log.Println(strings.TrimSpace(string(out)))
	if err != nil {
		log.Fatal(err)
	}
}
