// Package cache implements sqlc's on-disk cache, designed after Bazel's
// local disk cache.
//
// It has two halves:
//
//   - A content-addressable store (CAS) holding blobs keyed by the SHA-256
//     hash of their contents, laid out as cas/<xx>/<hash>.
//   - An action cache (AC) mapping the digest of an action — a description of
//     cacheable work and all of its inputs — to the digests of the outputs
//     that work produced, laid out as ac/<xx>/<hash>.
//
// Work whose output is not derivable from its inputs, like query analysis,
// uses both halves, just like Bazel: hash the action, look its digest up in
// the action cache, then fetch the referenced output blobs from the CAS.
//
// Remote fetches with a declared checksum, like WASM plugins, need no action
// cache entry at all: the declared sha256 is itself a content address, so the
// blob is stored and loaded directly from the CAS keyed by that checksum.
package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// Cache bundles the CAS and the action cache that shares it.
type Cache struct {
	root    string
	CAS     *CAS
	Actions *ActionCache
}

// Dir returns the cache root, defaulting to os.UserCacheDir(). The location
// can be overridden with the SQLCCACHE environment variable.
func Dir() (string, error) {
	cache := os.Getenv("SQLCCACHE")
	if cache != "" {
		return cache, nil
	}
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheHome, "sqlc"), nil
}

// Open returns the cache rooted at Dir().
func Open() (*Cache, error) {
	root, err := Dir()
	if err != nil {
		return nil, err
	}
	return OpenAt(root)
}

// OpenAt returns the cache rooted at the given directory, creating it if
// necessary.
func OpenAt(root string) (*Cache, error) {
	cas, err := newCAS(root)
	if err != nil {
		return nil, err
	}
	return &Cache{
		root:    root,
		CAS:     cas,
		Actions: newActionCache(root, cas),
	}, nil
}

// ExecDir returns a stable directory for materializing the output tree of
// the given action, for tools that need their outputs on disk (like wazero's
// compilation cache). It lives at exec/<action-hash> under the cache root
// and, like everything else in the cache, is safe to delete at any time: the
// authoritative copy of its contents is the CAS.
func (c *Cache) ExecDir(action Digest) (string, error) {
	dir := filepath.Join(c.root, "exec", action.Hash)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create %s directory: %w", dir, err)
	}
	return dir, nil
}
