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

	arch := flag.Arg(0)
	if arch == "" {
		log.Fatalf("missing platform_arch argument")
	}

	sha := os.Getenv("GITHUB_SHA")
	ref := os.Getenv("GITHUB_REF")
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

	xname := "./equinox"
	if _, err := os.Stat("./equinox"); os.IsNotExist(err) {
		xname = "equinox"
	}

	channel := "devel"
	if strings.HasPrefix(ref, "refs/tags/") {
		channel = "stable"
		version = strings.TrimPrefix(ref, "refs/tags/")
	}

	args := []string{"release",
		"--channel", channel,
		"--version", version,
	}

	if *draft {
		args = append(args, "--draft")
	}

	x := "-X github.com/kyleconroy/sqlc/internal/cmd.version=" + version
	args = append(args, []string{
		"--platforms", flag.Arg(0),
		"--app", "app_i4iCp1SuYfZ",
		"--token", os.Getenv("EQUINOX_API_TOKEN"),
		"--",
		"-ldflags", x, "./cmd/sqlc",
	}...)

	log.Printf("Releasing %s on channel %s", flag.Arg(0), channel)
	cmd = exec.Command(xname, args...)
	cmd.Env = os.Environ()
	out, err = cmd.CombinedOutput()
	log.Println(strings.TrimSpace(string(out)))
	if err != nil {
		log.Fatal(err)
	}
}
