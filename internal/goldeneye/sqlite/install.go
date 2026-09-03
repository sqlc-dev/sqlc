package sqlite

import (
	"archive/zip"
	"context"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// DefaultVersion is the SQLite release the dialect is generated from. It is
// the release the ncruces/go-sqlite3 driver in the main module embeds, so
// the functions the dialect knows are the ones the tests run against.
// Bumping it is a deliberate change: releases add functions and overloads,
// so regenerate and review the dialect after changing it, and add the new
// release's downloads to the table below.
const DefaultVersion = "3.53.4"

// asset is one downloadable build of the SQLite command-line tools: a zip
// published on sqlite.org holding the sqlite3 shell (sqlite3.exe on
// Windows) next to sqldiff, sqlite3_analyzer and sqlite3_rsync.
type asset struct {
	Version string
	OS      string
	Arch    string
	// Path is the download's address under https://sqlite.org/, the
	// release year included, as the download page's index lists it.
	Path string
	// SHA3 is the SHA3-256 of the zip, as the download page lists it.
	SHA3 string
}

// assets lists every build Install knows how to fetch, with the SHA3-256 of
// the download. A version that is not in this table cannot be installed:
// verifying the download is the point of the table.
//
// The checksums are the ones in the index at the foot of
// https://sqlite.org/download.html. sqlite.org builds the tools for x64
// Linux, both macOS architectures and both Windows architectures; there is
// no Linux arm64 build to list.
var assets = []asset{
	{"3.53.4", "linux", "amd64", "2026/sqlite-tools-linux-x64-3530400.zip", "6eeb57e8f2aef7687f9f016a980992cf2799c8c07a87c5e21495530f91915047"},
	{"3.53.4", "darwin", "amd64", "2026/sqlite-tools-osx-x64-3530400.zip", "3353bb4e5ac54f85c5b82012d30476d20539a9495abdd11a1707df578cff2d7e"},
	{"3.53.4", "darwin", "arm64", "2026/sqlite-tools-osx-arm64-3530400.zip", "58d53e0eb69c17cabebed2754bf399e4d44939be42dcf194769c00078bfd776d"},
	{"3.53.4", "windows", "amd64", "2026/sqlite-tools-win-x64-3530400.zip", "88b4659fe747896b853af10157316b4ade143553efb89c1c8ca7423a278dcc8b"},
	{"3.53.4", "windows", "arm64", "2026/sqlite-tools-win-arm64-3530400.zip", "0c99da3702b2517c1d738207db7e945e5c55be7748141a192a1c8f3b4455c44b"},
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
			return asset{}, fmt.Errorf("no SQLite %s build is listed for %s/%s", version, goos, goarch)
		}
	}
	return asset{}, fmt.Errorf("SQLite %s is not in the asset table; add its downloads and checksums to install.go", version)
}

// url is the asset's download address.
func (a asset) url() string {
	return "https://sqlite.org/" + a.Path
}

// binaryName is what the shell is called inside the zip and in the cache.
func binaryName(goos string) string {
	if goos == "windows" {
		return "sqlite3.exe"
	}
	return "sqlite3"
}

// cachedBinary is where Install puts the shell for a version.
func cachedBinary(version, goos string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sqlc-sqlite", version, binaryName(goos)), nil
}

// Locate finds a sqlite3 shell: the SQLITE3 environment variable wins, then
// the cached copy of DefaultVersion.
func Locate() (string, error) {
	if path := os.Getenv("SQLITE3"); path != "" {
		return path, nil
	}
	path, err := cachedBinary(DefaultVersion, runtime.GOOS)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("sqlite %s is not installed: run `go run ./cmd/goldeneye install sqlite` in internal/goldeneye, or set SQLITE3 to a sqlite3 shell", DefaultVersion)
	}
	return path, nil
}

// Install downloads the sqlite3 shell for a version into the cache and
// returns its path. It is a no-op when the version is already cached. The
// zip is checked against the table's SHA3-256 before the shell is taken
// out of it.
func Install(ctx context.Context, version, goos, goarch string, progress io.Writer) (string, error) {
	dest, err := cachedBinary(version, goos)
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
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}

	url := a.url()
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

	// The zip has to land on disk before it can be read — its directory is
	// at the end — so download it whole, next to the destination, and hash
	// every byte on the way.
	archive, err := os.CreateTemp(filepath.Dir(dest), "sqlite-tools-*.partial")
	if err != nil {
		return "", err
	}
	defer os.Remove(archive.Name())
	sum := sha3.New256()
	if _, err := io.Copy(io.MultiWriter(archive, sum), resp.Body); err != nil {
		archive.Close()
		return "", err
	}
	if err := archive.Close(); err != nil {
		return "", err
	}
	if got := hex.EncodeToString(sum.Sum(nil)); got != a.SHA3 {
		return "", fmt.Errorf("downloading %s: SHA3-256 mismatch: got %s, want %s", url, got, a.SHA3)
	}

	// Write the shell next to the destination and rename so a partial
	// extraction never masquerades as an installed binary.
	tmp, err := os.CreateTemp(filepath.Dir(dest), "sqlite3-*.partial")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	if err := extractBinary(archive.Name(), binaryName(goos), tmp); err != nil {
		tmp.Close()
		return "", fmt.Errorf("downloading %s: %w", url, err)
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

// extractBinary copies the named file out of a zip. The tools zips are
// flat, so the name is the whole path.
func extractBinary(archive, name string, dst io.Writer) error {
	zr, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		_, err = io.Copy(dst, rc)
		return err
	}
	return fmt.Errorf("zip does not contain %s", name)
}
