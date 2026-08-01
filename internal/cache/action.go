package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ActionResult records the outputs of a completed action, mirroring Bazel's
// ActionResult message. Outputs are not stored inline: each named output is a
// digest pointing into the CAS.
type ActionResult struct {
	// Outputs maps an output name (e.g. "analysis.pb", "plugin.wasm") to the
	// CAS digest of its contents.
	Outputs map[string]Digest `json:"outputs"`
}

// ActionCache maps action digests to ActionResults, stored as JSON files at
// ac/<xx>/<hash> under the cache root. Unlike CAS entries, action cache
// entries are not self-validating — the value is not derivable from the key —
// so Get additionally checks that every referenced output still exists in the
// CAS before reporting a hit, exactly like Bazel's disk cache does.
//
// Entry I/O goes through the same os.Root as the CAS, confining every read
// and write to the cache directory.
type ActionCache struct {
	root *os.Root
	cas  *CAS
}

func newActionCache(root *os.Root, cas *CAS) *ActionCache {
	return &ActionCache{root: root, cas: cas}
}

// path returns an entry's path relative to the cache root.
func (a *ActionCache) path(d Digest) string {
	return filepath.Join("ac", d.Hash[:2], d.Hash)
}

// Get returns the cached result for an action, or ErrNotFound on a miss. An
// entry whose outputs are missing or corrupt in the CAS is treated as a miss
// and evicted.
func (a *ActionCache) Get(action Digest) (*ActionResult, error) {
	if !action.valid() {
		return nil, ErrNotFound
	}
	path := a.path(action)
	data, err := a.root.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cache: %w", err)
	}
	var result ActionResult
	if err := json.Unmarshal(data, &result); err != nil {
		a.root.Remove(path)
		return nil, ErrNotFound
	}
	for _, d := range result.Outputs {
		if !a.cas.Contains(d) {
			a.root.Remove(path)
			return nil, ErrNotFound
		}
	}
	return &result, nil
}

// PutTree stores every file under dir in the CAS and records them as the
// action's outputs, named by their paths relative to dir. Use this for
// actions whose tool writes an output directory, like WASM compilation.
func (a *ActionCache) PutTree(action Digest, dir string) error {
	outputs := map[string]Digest{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		d, err := a.cas.Put(data)
		if err != nil {
			return err
		}
		outputs[filepath.ToSlash(rel)] = d
		return nil
	})
	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	if len(outputs) == 0 {
		return fmt.Errorf("cache: no outputs found under %s", dir)
	}
	return a.Put(action, &ActionResult{Outputs: outputs})
}

// GetTree materializes a cached action's outputs as files under dir, or
// returns ErrNotFound on a miss. A file already present is reused only if its
// contents still hash to the expected digest; otherwise it is rewritten from
// the CAS. Writes are staged, fsynced, and renamed, so a crash cannot leave a
// right-sized but torn file that later reads would trust.
func (a *ActionCache) GetTree(action Digest, dir string) error {
	result, err := a.Get(action)
	if err != nil {
		return err
	}
	for rel, d := range result.Outputs {
		// Output names come from the action cache entry, which — unlike a CAS
		// blob — is not self-validating, so a tampered entry could carry a
		// name like "../../etc/x". Reject anything that isn't a relative path
		// confined to dir; this is the confinement the os.Root gives the rest
		// of the cache, which GetTree can't use because dir is a caller-owned
		// path a tool must read by absolute name.
		if !filepath.IsLocal(rel) {
			return fmt.Errorf("cache: unsafe output path %q in action %s", rel, action)
		}
		path := filepath.Join(dir, filepath.FromSlash(rel))
		// Trust an existing file only if it still hashes to the digest;
		// size alone can mask in-place corruption that would make the
		// consuming tool hard-fail with no way to repair (see wazero).
		if existing, err := os.ReadFile(path); err == nil && DigestOf(existing) == d {
			continue
		}
		data, err := a.cas.Get(d)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
		if err := writeFileAtomic(path, data); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
	}
	return nil
}

// writeFileAtomic writes data to path via a staged temp file in the same
// directory, fsynced before an atomic rename, so a reader never observes a
// partial or torn file even across a crash.
func writeFileAtomic(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// Put records the result of an action. All outputs must already be in the
// CAS; writes are staged and renamed so concurrent processes never observe a
// partial entry.
func (a *ActionCache) Put(action Digest, result *ActionResult) error {
	for name, d := range result.Outputs {
		if !a.cas.Contains(d) {
			return fmt.Errorf("cache: output %q (%s) missing from CAS", name, d)
		}
	}
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	path := a.path(action)
	if err := a.root.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	f, name, err := a.cas.createTemp(action.Hash[:8] + "-")
	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	defer a.root.Remove(name)
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("cache: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	if err := a.root.Rename(name, path); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	return nil
}
