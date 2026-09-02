package clickhouse

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha512"
	"encoding/hex"
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
// after changing it, and add the new release's assets to the table below.
const DefaultVersion = "25.8.2.29"

// asset is one downloadable build of ClickHouse. Linux builds are published
// as clickhouse-common-static tarballs holding the binary at
// usr/bin/clickhouse; macOS builds are published as bare binaries.
type asset struct {
	Version string
	OS      string
	Arch    string
	Name    string
	SHA512  string
}

// assets lists every build Install knows how to fetch, with the SHA-512 of
// the download. A version that is not in this table cannot be installed:
// verifying the download is the point of the table.
//
// The tarball checksums are the ones in the .sha512 files ClickHouse
// publishes next to them. ClickHouse publishes no checksum for the macOS
// binaries, so those were computed from the downloads.
var assets = []asset{
	{"25.8.2.29", "linux", "amd64", "clickhouse-common-static-25.8.2.29-amd64.tgz", "6ff0aa1ffac6e564970174422ecde0d645cdb96812247a6e544d39cad6d78a514265f90a2bc7b4bad49903cea96eddd16a415a45b2aeaf9164461be76331bdee"},
	{"25.8.2.29", "linux", "arm64", "clickhouse-common-static-25.8.2.29-arm64.tgz", "68204ca4d4e472790f808ee376251fae82e58066a31f35a40d15d442ce5988d697f18a1208d28b8bb8e2dfad4b20b7fcb5107e2178472abcd97251b8de7f058e"},
	{"25.8.2.29", "darwin", "amd64", "clickhouse-macos", "2805805ad2506e37a3e71b4ae9e797bdc010a9368dc28e99bcaaa2c70a72cfdd031c0fce8fc304248fad73211d28d645f75c6432cfae1e1e54d72d04e8626cd4"},
	{"25.8.2.29", "darwin", "arm64", "clickhouse-macos-aarch64", "4c9237e85c8d4e1aced2b339b32e086b4f23fa14f99754b056a7e33dda88c4fb5d52e09e29d20181ec5b580caadb00f36613802e71755ff34927535eb8babf79"},
}

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

// releaseAsset finds the build for a platform in the table.
func releaseAsset(version, goos, goarch string) (asset, error) {
	for _, a := range assets {
		if a.Version == version && a.OS == goos && a.Arch == goarch {
			return a, nil
		}
	}
	for _, a := range assets {
		if a.Version == version {
			return asset{}, fmt.Errorf("no ClickHouse %s build is listed for %s/%s", version, goos, goarch)
		}
	}
	return asset{}, fmt.Errorf("ClickHouse %s is not in the asset table; add its downloads and checksums to install.go", version)
}

// url is the asset's download address on GitHub.
func (a asset) url() (string, error) {
	tag, err := releaseTag(a.Version)
	if err != nil {
		return "", err
	}
	return "https://github.com/ClickHouse/ClickHouse/releases/download/" + tag + "/" + a.Name, nil
}

// tarball reports whether the download is an archive rather than the binary.
func (a asset) tarball() bool {
	return strings.HasSuffix(a.Name, ".tgz")
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
		return "", fmt.Errorf("clickhouse %s is not installed: run `go run ./cmd/testcheck install clickhouse` in internal/testcheck, or set CLICKHOUSE to a clickhouse binary", DefaultVersion)
	}
	return path, nil
}

// Install downloads the clickhouse binary for a version into the cache and
// returns its path. It is a no-op when the version is already cached. The
// download is checked against the table's SHA-512 before it is installed.
func Install(ctx context.Context, version, goos, goarch string, progress io.Writer) (string, error) {
	dest, err := cachedBinary(version)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}
	a, err := releaseAsset(version, goos, goarch)
	if err != nil {
		return "", err
	}
	url, err := a.url()
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

	// Write next to the destination and rename so a partial or corrupt
	// download never masquerades as an installed binary.
	tmp, err := os.CreateTemp(filepath.Dir(dest), "clickhouse-*.partial")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	// Hash every byte that comes off the wire, including the parts of a
	// tarball after the binary, which extraction would otherwise not read.
	sum := sha512.New()
	body := io.TeeReader(resp.Body, sum)
	var src io.Reader = body
	if a.tarball() {
		src, err = binaryInTarball(body)
		if err != nil {
			return "", fmt.Errorf("downloading %s: %w", url, err)
		}
	}
	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		return "", err
	}
	if _, err := io.Copy(io.Discard, body); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if got := hex.EncodeToString(sum.Sum(nil)); got != a.SHA512 {
		return "", fmt.Errorf("downloading %s: SHA-512 mismatch: got %s, want %s", url, got, a.SHA512)
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
