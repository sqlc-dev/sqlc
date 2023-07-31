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
	docker := flag.Bool("docker", false, "create a docker release")
	flag.Parse()

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

	if *docker {
		x := "-extldflags \"-static\" -X github.com/sqlc-dev/sqlc/internal/cmd.version=" + version
		args := []string{
			"build",
			"-a",
			"-ldflags", x,
			"-o", "/workspace/sqlc",
			"./cmd/sqlc",
		}
		cmd := exec.Command("go", args...)
		cmd.Env = os.Environ()
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(strings.TrimSpace(string(out)))
			log.Fatal(err)
		}
		return
	}

	arch := flag.Arg(0)
	if arch == "" {
		log.Fatalf("missing platform_arch argument")
	}

	log.Fatal("publishing to Equinox has been disabled")
}
