package main

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
	"strconv"
	"strings"
)

// DefaultVersion is the ClickHouse release the goldens are generated with.
// Bumping it is a deliberate change: the query tree format and type
// inference can shift between releases, so regenerate and review the goldens
// after changing it.
const DefaultVersion = "25.8.2.29"

// releaseTag returns the GitHub release tag for a version. ClickHouse tags
// its March and August releases as LTS and everything else as stable.
func releaseTag(version string) (string, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid ClickHouse version %q: want MAJOR.MINOR.PATCH.BUILD", version)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid ClickHouse version %q: %w", version, err)
	}
	suffix := "-stable"
	if minor == 3 || minor == 8 {
		suffix = "-lts"
	}
	return "v" + version + suffix, nil
}

// releaseAsset returns the download URL for a platform and whether it is a
// tarball holding the binary at usr/bin/clickhouse rather than the bare
// binary. Linux builds are only published as tarballs; macOS builds only as
// bare binaries.
func releaseAsset(version, goos, goarch string) (url string, tarball bool, err error) {
	tag, err := releaseTag(version)
	if err != nil {
		return "", false, err
	}
	base := "https://github.com/ClickHouse/ClickHouse/releases/download/" + tag + "/"
	switch goos + "/" + goarch {
	case "linux/amd64":
		return base + "clickhouse-common-static-" + version + "-amd64.tgz", true, nil
	case "linux/arm64":
		return base + "clickhouse-common-static-" + version + "-arm64.tgz", true, nil
	case "darwin/amd64":
		return base + "clickhouse-macos", false, nil
	case "darwin/arm64":
		return base + "clickhouse-macos-aarch64", false, nil
	}
	return "", false, fmt.Errorf("no ClickHouse build is published for %s/%s", goos, goarch)
}

// cachedBinary is where Install puts the binary for a version.
func cachedBinary(version string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sqlc-clickhouse", version, "clickhouse"), nil
}

// Locate finds a clickhouse binary: the CLICKHOUSE environment variable wins,
// then the cached copy of DefaultVersion.
func Locate() (string, error) {
	if path := os.Getenv("CLICKHOUSE"); path != "" {
		return path, nil
	}
	path, err := cachedBinary(DefaultVersion)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("clickhouse %s is not installed: run `testgen install` or set CLICKHOUSE to a clickhouse binary", DefaultVersion)
	}
	return path, nil
}

// Install downloads the clickhouse binary for a version into the cache and
// returns its path. It is a no-op when the version is already cached.
func Install(ctx context.Context, version, goos, goarch string, progress io.Writer) (string, error) {
	dest, err := cachedBinary(version)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}
	url, tarball, err := releaseAsset(version, goos, goarch)
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
	tmp, err := os.CreateTemp(filepath.Dir(dest), "clickhouse-*.partial")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	var src io.Reader = resp.Body
	if tarball {
		src, err = binaryInTarball(resp.Body)
		if err != nil {
			return "", fmt.Errorf("downloading %s: %w", url, err)
		}
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

// binaryInTarball positions a reader at the clickhouse binary inside a
// clickhouse-common-static tarball.
func binaryInTarball(r io.Reader) (io.Reader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, errors.New("tarball does not contain usr/bin/clickhouse")
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag == tar.TypeReg && strings.HasSuffix(hdr.Name, "/usr/bin/clickhouse") {
			return tr, nil
		}
	}
}
