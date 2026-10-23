package duckdb

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DefaultVersion is the DuckDB release the dialect is generated from, the
// release darkwing is pinned against. Bumping it is a deliberate change:
// every release adds functions and types, so regenerate and review the
// dialect after changing it, and add the new release's assets to the table
// below.
const DefaultVersion = "2.0.0"

// asset is one downloadable build of the DuckDB CLI, a tarball holding the
// duckdb binary at its top level. DuckDB publishes one macOS build for both
// architectures.
type asset struct {
	Version string
	OS      string
	Arch    string
	Name    string
	SHA256  string
}

// assets lists every build Install knows how to fetch, with the SHA-256 of
// the download. A version that is not in this table cannot be installed:
// verifying the download is the point of the table.
//
// DuckDB publishes no checksum files next to its release assets, so these
// were computed from the downloads.
var assets = []asset{
	{"2.0.0", "linux", "amd64", "duckdb-cli-linux-amd64.tar.gz", "c6507206baf7793c6013dbd82368469a45f20dc0e17a6f92f204602f652ed2bc"},
	{"2.0.0", "linux", "arm64", "duckdb-cli-linux-arm64.tar.gz", "859e1875d2a9515fbb4bf8871061437a14af82320da68c753f094c4b776348fe"},
	{"2.0.0", "darwin", "amd64", "duckdb-cli-osx-universal.tar.gz", "8a8f67e7c7f7d1392fa63ba3971948391392fe2a25fca63e7008e33c50beb4c2"},
	{"2.0.0", "darwin", "arm64", "duckdb-cli-osx-universal.tar.gz", "8a8f67e7c7f7d1392fa63ba3971948391392fe2a25fca63e7008e33c50beb4c2"},
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
			return asset{}, fmt.Errorf("no DuckDB %s build is listed for %s/%s", version, goos, goarch)
		}
	}
	return asset{}, fmt.Errorf("DuckDB %s is not in the asset table; add its downloads and checksums to install.go", version)
}

// url is the asset's download address on GitHub.
func (a asset) url() string {
	return "https://github.com/duckdb/duckdb/releases/download/v" + a.Version + "/" + a.Name
}

// cachedBinary is where Install puts the binary for a version.
func cachedBinary(version string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sqlc-duckdb", version, "duckdb"), nil
}

// Install downloads the CLI for a version into the cache and returns its
// path. It is a no-op when the version is already cached. The download is
// checked against the table's SHA-256 before it is installed.
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
	url := a.url()
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
	tmp, err := os.CreateTemp(filepath.Dir(dest), "duckdb-*.partial")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	// Hash every byte that comes off the wire, including the parts of the
	// tarball after the binary, which extraction would otherwise not read.
	sum := sha256.New()
	body := io.TeeReader(resp.Body, sum)
	src, err := binaryInTarball(body)
	if err != nil {
		tmp.Close()
		return "", fmt.Errorf("downloading %s: %w", url, err)
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
	if got := hex.EncodeToString(sum.Sum(nil)); got != a.SHA256 {
		return "", fmt.Errorf("downloading %s: SHA-256 mismatch: got %s, want %s", url, got, a.SHA256)
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
