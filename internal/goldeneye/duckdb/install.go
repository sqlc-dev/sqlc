package duckdb

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DefaultVersion is the DuckDB preview channel the CLI is downloaded from:
// the v2.0 development builds, which DuckDB publishes as a rolling tarball
// per platform under https://artifacts.duckdb.org/<channel>/ until 2.0 is
// released. A channel has no per-build download and no checksum to pin,
// so Install fetches whatever build the channel holds and the version the
// CLI reports is what a run logs; GeneratedFrom records the build the
// committed dialect came from, and a check against a later build reports
// what the later build added. Once 2.0 is released, pin the release and
// its checksums here the way the clickhouse package does.
const DefaultVersion = "v2.0-cyanoptera"

// GeneratedFrom is the build the committed dialect was generated from, as
// `duckdb --version` reports it. Update it when regenerating.
const GeneratedFrom = "v2.0.0-alpha41396 (Cyanoptera) d41e527b18"

// channelURL is the download address of a channel's CLI tarball for a
// platform. DuckDB publishes one macOS build for both architectures.
func channelURL(channel, goos, goarch string) (string, error) {
	var name string
	switch {
	case goos == "linux" && (goarch == "amd64" || goarch == "arm64"):
		name = "duckdb-cli-linux-" + goarch + ".tar.gz"
	case goos == "darwin":
		name = "duckdb-cli-osx-universal.tar.gz"
	default:
		return "", fmt.Errorf("no DuckDB %s build is published for %s/%s", channel, goos, goarch)
	}
	return "https://artifacts.duckdb.org/" + channel + "/" + name, nil
}

// cachedBinary is where Install puts the binary for a channel.
func cachedBinary(channel string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sqlc-duckdb", channel, "duckdb"), nil
}

// Install downloads the channel's current CLI into the cache and returns
// its path. It is a no-op when the channel is already cached: remove the
// cached directory to fetch the channel's newer build.
func Install(ctx context.Context, channel, goos, goarch string, progress io.Writer) (string, error) {
	dest, err := cachedBinary(channel)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}
	url, err := channelURL(channel, goos, goarch)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}

	fmt.Fprintf(progress, "downloading %s\n", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading %s: %s", url, resp.Status)
	}

	// Write next to the destination and rename so a partial download never
	// masquerades as an installed binary.
	tmp, err := os.CreateTemp(filepath.Dir(dest), "duckdb-*.partial")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	src, err := binaryInTarball(resp.Body)
	if err != nil {
		tmp.Close()
		return "", fmt.Errorf("downloading %s: %w", url, err)
	}
	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(tmp.Name(), dest); err != nil {
		return "", err
	}
	return dest, nil
}

// binaryInTarball positions a reader at the duckdb binary inside a CLI
// tarball, which holds it at the top level.
func binaryInTarball(r io.Reader) (io.Reader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, errors.New("tarball does not contain duckdb")
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg && strings.TrimPrefix(hdr.Name, "./") == "duckdb" {
			return tr, nil
		}
	}
}
