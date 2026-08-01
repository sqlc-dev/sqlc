package cache

import (
	"encoding/json"
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
type ActionCache struct {
	root string
	cas  *CAS
}

func newActionCache(root string, cas *CAS) *ActionCache {
	return &ActionCache{root: root, cas: cas}
}

func (a *ActionCache) path(d Digest) string {
	return filepath.Join(a.root, "ac", d.Hash[:2], d.Hash)
}

// Get returns the cached result for an action, or ErrNotFound on a miss. An
// entry whose outputs are missing or corrupt in the CAS is treated as a miss
// and evicted.
func (a *ActionCache) Get(action Digest) (*ActionResult, error) {
	if !action.valid() {
		return nil, ErrNotFound
	}
	path := a.path(action)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("cache: %w", err)
	}
	var result ActionResult
	if err := json.Unmarshal(data, &result); err != nil {
		os.Remove(path)
		return nil, ErrNotFound
	}
	for _, d := range result.Outputs {
		if !a.cas.Contains(d) {
			os.Remove(path)
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
// returns ErrNotFound on a miss. Files already present with the right size
// are left in place; missing ones are staged and renamed so concurrent
// processes never observe partial files.
func (a *ActionCache) GetTree(action Digest, dir string) error {
	result, err := a.Get(action)
	if err != nil {
		return err
	}
	for rel, d := range result.Outputs {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if fi, err := os.Stat(path); err == nil && fi.Size() == d.SizeBytes {
			continue
		}
		data, err := a.cas.Get(d)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("cache: %w", err)
		}
		f, err := os.CreateTemp(a.cas.tmp, d.Hash[:8]+"-*")
		if err != nil {
			return fmt.Errorf("cache: %w", err)
		}
		if _, err := f.Write(data); err != nil {
			f.Close()
			os.Remove(f.Name())
			return fmt.Errorf("cache: %w", err)
		}
		if err := f.Close(); err != nil {
			os.Remove(f.Name())
			return fmt.Errorf("cache: %w", err)
		}
		if err := os.Rename(f.Name(), path); err != nil {
			os.Remove(f.Name())
			return fmt.Errorf("cache: %w", err)
		}
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
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	f, err := os.CreateTemp(a.cas.tmp, action.Hash[:8]+"-*")
	if err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return fmt.Errorf("cache: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	return nil
}
