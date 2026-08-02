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

// Cache bundles the CAS and the action cache that shares it. All storage
// I/O is confined to the cache directory through an os.Root; callers should
// Close the cache when finished with it to release the root.
type Cache struct {
	root    *os.Root
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
func OpenAt(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create %s directory: %w", dir, err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}
	cas, err := newCAS(root)
	if err != nil {
		root.Close()
		return nil, err
	}
	return &Cache{
		root:    root,
		CAS:     cas,
		Actions: newActionCache(root, cas),
	}, nil
}

// Close releases the cache's handle on its root directory.
func (c *Cache) Close() error {
	return c.root.Close()
}

// ExecDir creates a fresh, private scratch directory for materializing the
// output tree of the given action, for tools that need their outputs on disk
// (like wazero's compilation cache). Each call returns a new directory under
// exec/, so two concurrent processes never share one — otherwise a tool
// staging files there (wazero writes <key>.tmp files in place) could be swept
// into the other's PutTree. The caller must remove it when done; its contents
// are always reproducible from the CAS, so losing it is harmless.
func (c *Cache) ExecDir(action Digest) (string, error) {
	base := filepath.Join(c.root.Name(), "exec")
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", fmt.Errorf("failed to create %s directory: %w", base, err)
	}
	dir, err := os.MkdirTemp(base, action.Hash+"-")
	if err != nil {
		return "", fmt.Errorf("cache: %w", err)
	}
	return dir, nil
}
