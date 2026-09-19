package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// clickhouseVersion is the ClickHouse release to install. It is the
	// release internal/goldeneye generates the dialect from, and the binary
	// is cached where goldeneye caches it, so the two share one download.
	clickhouseVersion = "25.8.2.29"

	// clickhousePassword is the password the default user is given: the
	// one docker-compose.yml and CI use.
	clickhousePassword = "mysecretpassword"

	clickhouseTCPPort  = "9000"
	clickhouseHTTPPort = "8123"
)

// clickhouseAsset is one downloadable build of ClickHouse: a
// clickhouse-common-static tarball holding the binary at usr/bin/clickhouse,
// with the SHA-512 ClickHouse publishes next to it.
type clickhouseAsset struct {
	URL    string
	SHA512 string
}

var clickhouseAssets = map[string]clickhouseAsset{
	"linux/amd64": {
		URL:    "https://github.com/ClickHouse/ClickHouse/releases/download/v" + clickhouseVersion + "-lts/clickhouse-common-static-" + clickhouseVersion + "-amd64.tgz",
		SHA512: "6ff0aa1ffac6e564970174422ecde0d645cdb96812247a6e544d39cad6d78a514265f90a2bc7b4bad49903cea96eddd16a415a45b2aeaf9164461be76331bdee",
	},
	"linux/arm64": {
		URL:    "https://github.com/ClickHouse/ClickHouse/releases/download/v" + clickhouseVersion + "-lts/clickhouse-common-static-" + clickhouseVersion + "-arm64.tgz",
		SHA512: "68204ca4d4e472790f808ee376251fae82e58066a31f35a40d15d442ce5988d697f18a1208d28b8bb8e2dfad4b20b7fcb5107e2178472abcd97251b8de7f058e",
	},
}

// clickhouseDir is where the release is cached, and the server's
// configuration and data are kept under it.
func clickhouseDir() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, "sqlc-clickhouse", clickhouseVersion), nil
}

func clickhouseBinary(dir string) string {
	return filepath.Join(dir, "clickhouse")
}

// installClickHouse downloads the ClickHouse release for this platform and
// unpacks its binary into the cache, checking the download against the
// pinned SHA-512. It is a no-op when the binary is already there, and
// skips platforms the release is not published for.
func installClickHouse() error {
	log.Printf("--- Installing ClickHouse %s ---", clickhouseVersion)

	platform := runtime.GOOS + "/" + runtime.GOARCH
	asset, ok := clickhouseAssets[platform]
	if !ok {
		log.Printf("clickhouse is not published for %s, skipping", platform)
		return nil
	}

	dir, err := clickhouseDir()
	if err != nil {
		return err
	}
	binary := clickhouseBinary(dir)
	if _, err := os.Stat(binary); err == nil {
		log.Printf("clickhouse %s is already installed at %s", clickhouseVersion, binary)
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	archive := filepath.Join(dir, "clickhouse.tgz")
	log.Printf("downloading %s", asset.URL)
	if err := downloadFile(archive, asset.URL); err != nil {
		return fmt.Errorf("downloading clickhouse: %w", err)
	}
	defer os.Remove(archive)

	sum, err := sha512File(archive)
	if err != nil {
		return err
	}
	if sum != asset.SHA512 {
		return fmt.Errorf("clickhouse download has SHA-512 %s, want %s", sum, asset.SHA512)
	}

	log.Printf("unpacking the binary into %s", binary)
	if err := extractClickHouse(archive, binary); err != nil {
		return fmt.Errorf("unpacking clickhouse: %w", err)
	}
	return nil
}

// extractClickHouse copies usr/bin/clickhouse out of the release tarball.
func extractClickHouse(archive, binary string) error {
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return errors.New("tarball does not contain usr/bin/clickhouse")
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg || !strings.HasSuffix(hdr.Name, "/usr/bin/clickhouse") {
			continue
		}
		tmp, err := os.CreateTemp(filepath.Dir(binary), "clickhouse-*.partial")
		if err != nil {
			return err
		}
		if _, err := io.Copy(tmp, tr); err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
			return err
		}
		if err := tmp.Close(); err != nil {
			os.Remove(tmp.Name())
			return err
		}
		if err := os.Chmod(tmp.Name(), 0o755); err != nil {
			os.Remove(tmp.Name())
			return err
		}
		return os.Rename(tmp.Name(), binary)
	}
}

// sha512File computes the SHA-512 hash of a file and returns the hex string.
func sha512File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// clickhouseConfig is the server's configuration: it listens on localhost
// only, keeps its data next to the configuration, and takes its users from
// users.xml beside it.
const clickhouseConfig = `<clickhouse>
  <logger>
    <level>warning</level>
    <console>1</console>
  </logger>
  <path>./data/</path>
  <tcp_port>` + clickhouseTCPPort + `</tcp_port>
  <http_port>` + clickhouseHTTPPort + `</http_port>
  <listen_host>127.0.0.1</listen_host>
  <users_config>users.xml</users_config>
  <user_directories>
    <users_xml>
      <path>users.xml</path>
    </users_xml>
  </user_directories>
</clickhouse>
`

// clickhouseUsers gives the default user a password. Without one the
// server restricts the user to localhost, which is fine here, but the
// tests, docker-compose.yml and CI all connect with the same password.
const clickhouseUsers = `<clickhouse>
  <profiles>
    <default/>
  </profiles>
  <quotas>
    <default/>
  </quotas>
  <users>
    <default>
      <password>` + clickhousePassword + `</password>
      <networks>
        <ip>::/0</ip>
      </networks>
      <profile>default</profile>
      <quota>default</quota>
    </default>
  </users>
</clickhouse>
`

// startClickHouse starts a ClickHouse server in the background, serving
// on localhost, and waits until it answers a query. It is a no-op when a
// server already does, and skips platforms the release is not installed
// on.
func startClickHouse() error {
	log.Println("--- Starting ClickHouse ---")

	dir, err := clickhouseDir()
	if err != nil {
		return err
	}
	binary := clickhouseBinary(dir)
	if _, err := os.Stat(binary); err != nil {
		if _, ok := clickhouseAssets[runtime.GOOS+"/"+runtime.GOARCH]; !ok {
			log.Printf("clickhouse is not published for %s/%s, skipping", runtime.GOOS, runtime.GOARCH)
			return nil
		}
		return fmt.Errorf("clickhouse is not installed: run `sqlc-test-setup install clickhouse` first")
	}

	if clickhouseReady(binary) {
		log.Println("clickhouse is already running and accepting connections")
		return nil
	}

	server := filepath.Join(dir, "server")
	if err := os.MkdirAll(filepath.Join(server, "data"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(server, "config.xml"), []byte(clickhouseConfig), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(server, "users.xml"), []byte(clickhouseUsers), 0o644); err != nil {
		return err
	}
	logFile, err := os.Create(filepath.Join(server, "server.log"))
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmd := exec.Command(binary, "server", "--config-file="+filepath.Join(server, "config.xml"))
	cmd.Dir = server
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = detachedProcess()
	log.Printf("starting %s server in %s", binary, server)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting clickhouse: %w", err)
	}
	if err := cmd.Process.Release(); err != nil {
		return err
	}

	log.Println("waiting for clickhouse to accept connections")
	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		if clickhouseReady(binary) {
			log.Println("clickhouse is accepting connections")
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("clickhouse did not start in time (see %s)", logFile.Name())
}

// clickhouseReady reports whether a server answers on the TCP port with
// the expected password.
func clickhouseReady(binary string) bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+clickhouseTCPPort, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return exec.Command(binary, "client", "--host", "127.0.0.1", "--port", clickhouseTCPPort,
		"--password", clickhousePassword, "-q", "SELECT 1").Run() == nil
}
