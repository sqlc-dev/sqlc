package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNotFound is returned when a blob or action result is not in the cache.
var ErrNotFound = errors.New("cache: not found")

// CAS is an on-disk content-addressable store modeled on Bazel's disk cache.
// Blobs live at cas/<function>/<xx>/<hash> under the cache root, where <xx>
// is the first two hex characters of the hash. Because a blob's name is
// derived from its contents, entries never change once written: writers race
// benignly and readers can detect corruption by re-hashing.
//
// Blobs are keyed by BLAKE3 by default. Content with an externally declared
// SHA-256 checksum (remotely fetched plugins) is stored keyed by that
// checksum instead — see PutSHA256 — so fetches need no action cache entry:
// the declared checksum is the address.
type CAS struct {
	root string
	tmp  string
}

func newCAS(root string) (*CAS, error) {
	tmp := filepath.Join(root, "tmp")
	if err := os.MkdirAll(tmp, 0755); err != nil {
		return nil, fmt.Errorf("cache: create %s: %w", tmp, err)
	}
	return &CAS{root: root, tmp: tmp}, nil
}

func (c *CAS) path(d Digest) string {
	return filepath.Join(c.root, "cas", string(d.function()), d.Hash[:2], d.Hash)
}

// Put stores a blob keyed by its BLAKE3 hash and returns its digest.
func (c *CAS) Put(data []byte) (Digest, error) {
	return c.put(DigestOf(data), data)
}

// PutSHA256 stores a blob keyed by its SHA-256 hash and returns its digest.
// Use this for content addressed by a checksum declared outside the cache,
// so later runs can load it with a bare SHA256Digest lookup.
func (c *CAS) PutSHA256(data []byte) (Digest, error) {
	hash, _ := sum(SHA256, data)
	return c.put(Digest{Function: SHA256, Hash: hash, SizeBytes: int64(len(data))}, data)
}

// put writes a blob atomically: it is staged in a temp file and renamed into
// place, so concurrent sqlc processes never observe partial entries.
func (c *CAS) put(d Digest, data []byte) (Digest, error) {
	path := c.path(d)
	// Skip the write only when an entry of the right size already exists; a
	// wrong-sized entry is corrupt and is atomically replaced by the rename
	// below.
	if fi, err := os.Stat(path); err == nil && fi.Size() == d.SizeBytes {
		return d, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	f, err := os.CreateTemp(c.tmp, d.Hash[:8]+"-*")
	if err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	if err := f.Close(); err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return Digest{}, fmt.Errorf("cache: %w", err)
	}
	return d, nil
}

// Get returns the blob for a digest. Contents are re-hashed with the
// digest's hash function before being returned; a corrupt entry is evicted
// and reported as ErrNotFound so callers simply redo the work that produced
// it.
func (c *CAS) Get(d Digest) ([]byte, error) {
	if !d.valid() {
		return nil, ErrNotFound
	}
	path := c.path(d)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cache: %w", err)
	}
	if hash, _ := sum(d.function(), data); hash != d.Hash {
		os.Remove(path)
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
	fi, err := os.Stat(c.path(d))
	if err != nil {
		return false
	}
	return d.SizeBytes < 0 || fi.Size() == d.SizeBytes
}
