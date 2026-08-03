package cache

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
)

// ErrNotFound is returned when a blob or action result is not in the cache.
var ErrNotFound = errors.New("cache: not found")

// CAS is an on-disk content-addressable store keyed by SHA-256, modeled on
// Bazel's disk cache. Blobs live at cas/<xx>/<hash> under the cache root,
// where <xx> is the first two hex characters of the hash. Because a blob's
// name is derived from its contents, entries never change once written:
// writers race benignly and readers can detect corruption by re-hashing.
//
// Content with an externally declared checksum (remotely fetched plugins)
// needs no action cache entry: the declared sha256 is the address, so it is
// stored and loaded directly — see SHA256Digest.
//
// All I/O goes through an os.Root, so no entry name — hash-derived or read
// from an action cache entry — can escape the cache directory.
type CAS struct {
	root *os.Root
}

func newCAS(root *os.Root) (*CAS, error) {
	if err := root.MkdirAll("tmp", 0755); err != nil {
		return nil, fmt.Errorf("cache: create tmp: %w", err)
	}
	return &CAS{root: root}, nil
}

// path returns a blob's path relative to the cache root.
func (c *CAS) path(d Digest) string {
	return filepath.Join("cas", d.Hash[:2], d.Hash)
}

// Filename returns the path of a stored blob, for a consumer that needs the
// file rather than its bytes — SQLite, for one, opens a database by name. A
// blob is named after the hash of its contents, so the file at this path never
// changes and any number of processes may read it at once.
//
// It reports false when the blob is not stored.
func (c *CAS) Filename(d Digest) (string, bool) {
	if !c.Contains(d) {
		return "", false
	}
	return filepath.Join(c.root.Name(), c.path(d)), true
}

// createTemp creates a staging file under tmp/ in the cache root, returning
// the open file and its root-relative name.
func (c *CAS) createTemp(prefix string) (*os.File, string, error) {
	for range 10000 {
		name := filepath.Join("tmp", prefix+strconv.FormatUint(rand.Uint64(), 36))
		f, err := c.root.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		return f, name, err
	}
	return nil, "", errors.New("cache: could not create temp file")
}

// Put stores a blob and returns its digest. Writing is atomic: the blob is
// staged in a temp file and renamed into place, so concurrent sqlc processes
// never observe partial entries.
func (c *CAS) Put(data []byte) (Digest, error) {
	d := DigestOf(data)
	path := c.path(d)
	// Skip the write only when an entry of the right size already exists; a
	// wrong-sized entry is corrupt and is atomically replaced by the rename
	// below.
	if fi, err := c.root.Stat(path); err == nil && fi.Size() == d.SizeBytes {
		return d, nil
	}
	if err := c.root.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	f, name, err := c.createTemp(d.Hash[:8] + "-")
	if err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	defer c.root.Remove(name)
	if _, err := f.Write(data); err != nil {
		f.Close()
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	if err := f.Close(); err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	if err := c.root.Rename(name, path); err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	return d, nil
}

// Get returns the blob for a digest. Contents are re-hashed before being
// returned; a corrupt entry is evicted and reported as ErrNotFound so
// callers simply redo the work that produced it.
func (c *CAS) Get(d Digest) ([]byte, error) {
	if !d.valid() {
		return nil, ErrNotFound
	}
	path := c.path(d)
	data, err := c.root.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cache: %w", err)
	}
	if DigestOf(data).Hash != d.Hash {
		c.root.Remove(path)
		return nil, ErrNotFound
	}
	return data, nil
}

// Contains reports whether a blob with the given digest is present, checking
// size when the digest carries one. It does not verify contents; Get performs
// full verification.
func (c *CAS) Contains(d Digest) bool {
	if !d.valid() {
		return false
	}
	fi, err := c.root.Stat(c.path(d))
	if err != nil {
		return false
	}
	return d.SizeBytes < 0 || fi.Size() == d.SizeBytes
}
