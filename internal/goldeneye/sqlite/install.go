package sqlite

import (
	"archive/zip"
	"context"
	"crypto/sha3"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// DefaultVersion is the SQLite release the dialect is generated from. It is
// the release the ncruces/go-sqlite3 driver in the main module embeds, so
// the functions the dialect knows are the ones the tests run against.
// Bumping it is a deliberate change: releases add functions and overloads,
// so regenerate and review the dialect after changing it, and add the new
// release's amalgamation to the table below.
const DefaultVersion = "3.53.4"

// defaultOptions are the compile options the base dialect is built with:
// the ones sqlite.org's own configure turns on by default, and so the ones
// a stock sqlite3 has. JSON is part of the library unless omitted, so only
// the math functions need naming.
var defaultOptions = []string{"SQLITE_ENABLE_MATH_FUNCTIONS"}

// extensions are the compile options that each become a directory under
// the dialect's extensions/, holding the functions a build with the option
// adds over the default build — the way each PostgreSQL contrib extension
// holds what CREATE EXTENSION adds. Each is named after its option as
// pragma compile_options spells it, in lower case, and lists the options
// its build needs: GEOPOLY lives inside the RTREE module, so enabling it
// alone adds nothing. Options that add virtual tables but no functions,
// such as SESSION and DBSTAT, are not listed, since the dialect has nothing
// to say about them, and ENABLE_FTS4 is the same module as ENABLE_FTS3.
var extensions = []build{
	{"soundex", []string{"SQLITE_SOUNDEX"}},
	{"enable_fts3", []string{"SQLITE_ENABLE_FTS3"}},
	{"enable_fts5", []string{"SQLITE_ENABLE_FTS5"}},
	{"enable_geopoly", []string{"SQLITE_ENABLE_RTREE", "SQLITE_ENABLE_GEOPOLY"}},
	{"enable_offset_sql_func", []string{"SQLITE_ENABLE_OFFSET_SQL_FUNC"}},
	{"enable_percentile", []string{"SQLITE_ENABLE_PERCENTILE"}},
	{"enable_rtree", []string{"SQLITE_ENABLE_RTREE"}},
}

// asset is one downloadable amalgamation of SQLite: the zip published on
// sqlite.org holding sqlite3.c, sqlite3.h and the shell's shell.c.
type asset struct {
	Version string
	// Path is the download's address under https://sqlite.org/, the
	// release year included, as the download page's index lists it.
	Path string
	// SHA3 is the SHA3-256 of the zip, as the download page lists it.
	SHA3 string
}

// assets lists every amalgamation Install knows how to fetch, with the
// SHA3-256 of the download. A version that is not in this table cannot be
// installed: verifying the download is the point of the table. The
// checksums are the ones in the index at the foot of
// https://sqlite.org/download.html.
var assets = []asset{
	{"3.53.4", "2026/sqlite-amalgamation-3530400.zip", "628a44cfe82c66aed1ccbbe85a562d2e33ebe64b3288981ed76285612227934e"},
}

// sources are the files taken out of the amalgamation.
var sources = []string{"sqlite3.c", "sqlite3.h", "shell.c"}

func releaseAsset(version string) (asset, error) {
	for _, a := range assets {
		if a.Version == version {
			return a, nil
		}
	}
	return asset{}, fmt.Errorf("SQLite %s is not in the asset table; add its amalgamation and checksum to install.go", version)
}

func (a asset) url() string {
	return "https://sqlite.org/" + a.Path
}

// cacheDir is where Install puts a version: the sources under src/, and one
// shell per build under default/ and under each option's extension name.
func cacheDir(version string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sqlc-sqlite", version), nil
}

// build is one compiled shell: its name, which is the directory it is
// under in the version's cache directory and, for an extension, under the
// dialect's extensions/, and the options it is built with beyond the
// default ones.
type build struct {
	name    string
	options []string
}

// builds lists the default build first, then one per extension.
func builds() []build {
	return append([]build{{"default", nil}}, extensions...)
}

// flags are every option a build is compiled with.
func (b build) flags() []string {
	return append(append([]string{}, defaultOptions...), b.options...)
}

func (b build) binary(dir string) string {
	return filepath.Join(dir, b.name, "sqlite3")
}

// Locate finds the directory of shells Install made for DefaultVersion.
func Locate() (string, error) {
	dir, err := cacheDir(DefaultVersion)
	if err != nil {
		return "", err
	}
	for _, b := range builds() {
		if _, err := os.Stat(b.binary(dir)); err != nil {
			return "", fmt.Errorf("sqlite %s is not built with %s: run `go run ./cmd/goldeneye install sqlite` in internal/goldeneye", DefaultVersion, strings.Join(b.flags(), " "))
		}
	}
	return dir, nil
}

// Install downloads the amalgamation for a version into the cache and
// compiles the shell from it once per build, returning the directory the
// shells are under. Builds already there are kept, so adding an option
// compiles only its shell. The download is checked against the table's
// SHA3-256 before it is unpacked. Compiling takes the compiler CC names, or
// cc, and a few seconds per build without optimisation, which a shell that
// only reads catalogs does not need.
//
// goos and goarch are what every installer is handed; the sources build the
// same everywhere a C compiler is, but no default compiler or link line is
// known for Windows.
func Install(ctx context.Context, version, goos, goarch string, progress io.Writer) (string, error) {
	if goos == "windows" {
		return "", errors.New("building SQLite from source is not supported on Windows")
	}
	dir, err := cacheDir(version)
	if err != nil {
		return "", err
	}
	var missing []build
	for _, b := range builds() {
		if _, err := os.Stat(b.binary(dir)); err != nil {
			missing = append(missing, b)
		}
	}
	if len(missing) == 0 {
		return dir, nil
	}
	if err := download(ctx, version, dir, progress); err != nil {
		return "", err
	}
	if err := compile(ctx, dir, missing, progress); err != nil {
		return "", err
	}
	return dir, nil
}

// download fetches the amalgamation and unpacks the sources into src/,
// unless they are already there.
func download(ctx context.Context, version, dir string, progress io.Writer) error {
	src := filepath.Join(dir, "src")
	have := true
	for _, name := range sources {
		if _, err := os.Stat(filepath.Join(src, name)); err != nil {
			have = false
		}
	}
	if have {
		return nil
	}
	a, err := releaseAsset(version)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(src, 0o755); err != nil {
		return err
	}

	url := a.url()
	fmt.Fprintf(progress, "downloading %s\n", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading %s: %s", url, resp.Status)
	}

	// The zip has to land on disk before it can be read — its directory is
	// at the end — so download it whole and hash every byte on the way.
	archive, err := os.CreateTemp(dir, "sqlite-amalgamation-*.partial")
	if err != nil {
		return err
	}
	defer os.Remove(archive.Name())
	sum := sha3.New256()
	if _, err := io.Copy(io.MultiWriter(archive, sum), resp.Body); err != nil {
		archive.Close()
		return err
	}
	if err := archive.Close(); err != nil {
		return err
	}
	if got := hex.EncodeToString(sum.Sum(nil)); got != a.SHA3 {
		return fmt.Errorf("downloading %s: SHA3-256 mismatch: got %s, want %s", url, got, a.SHA3)
	}

	zr, err := zip.OpenReader(archive.Name())
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, name := range sources {
		if err := extract(zr, name, filepath.Join(src, name)); err != nil {
			return fmt.Errorf("downloading %s: %w", url, err)
		}
	}
	return nil
}

// extract copies the named file out of the zip. The amalgamation's files
// sit in one directory named after the release, so the name is matched on
// its base.
func extract(zr *zip.ReadCloser, name, dest string) error {
	for _, f := range zr.File {
		if filepath.Base(f.Name) != name || f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		tmp, err := os.CreateTemp(filepath.Dir(dest), name+"-*.partial")
		if err != nil {
			return err
		}
		defer os.Remove(tmp.Name())
		if _, err := io.Copy(tmp, rc); err != nil {
			tmp.Close()
			return err
		}
		if err := tmp.Close(); err != nil {
			return err
		}
		return os.Rename(tmp.Name(), dest)
	}
	return fmt.Errorf("zip does not contain %s", name)
}

// compile builds the shells, as many at a time as there are CPUs.
func compile(ctx context.Context, dir string, todo []build, progress io.Writer) error {
	cc := os.Getenv("CC")
	if cc == "" {
		cc = "cc"
	}
	if _, err := exec.LookPath(cc); err != nil {
		return fmt.Errorf("no C compiler found: put cc on PATH or set CC to one")
	}
	var wg sync.WaitGroup
	slots := make(chan struct{}, max(1, runtime.NumCPU()))
	errs := make([]error, len(todo))
	for i, b := range todo {
		wg.Add(1)
		go func() {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			errs[i] = b.compile(ctx, cc, dir, progress)
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

// compile builds one shell from the sources in src/, writing it beside its
// destination and renaming so a failed build never masquerades as one.
func (b build) compile(ctx context.Context, cc, dir string, progress io.Writer) error {
	dest := b.binary(dir)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp := dest + ".partial"
	defer os.Remove(tmp)
	args := []string{"-O0"}
	for _, opt := range b.flags() {
		args = append(args, "-D"+opt)
	}
	src := filepath.Join(dir, "src")
	args = append(args, filepath.Join(src, "shell.c"), filepath.Join(src, "sqlite3.c"), "-o", tmp, "-lm", "-ldl", "-lpthread")
	fmt.Fprintf(progress, "building %s: %s %s\n", b.name, cc, strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, cc, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("building %s: %w\n%s", b.name, err, strings.TrimSpace(string(out)))
	}
	return os.Rename(tmp, dest)
}
