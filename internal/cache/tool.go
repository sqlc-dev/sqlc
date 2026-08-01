package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/info"
)

// The sha256 of the sqlc binary is an input to every action — a rebuilt sqlc
// may analyze queries or embed a different wazero than the one that produced
// a cache entry, even when the version string is unchanged (dev builds).
//
// Hashing a ~100MB executable costs tens of milliseconds, too much to pay on
// every short-lived sqlc process, so the digest is memoized on disk keyed by
// the executable's path, size, and mtime — the same trick Bazel's file
// digest cache uses. A warm run costs one stat and a tiny read; only a
// rebuilt (or moved) binary is re-hashed.
var tool struct {
	sync.Mutex
	digest []byte
}

type toolMemo struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	MtimeNS   int64  `json:"mtime_ns"`
	SHA256    string `json:"sha256"`
}

const toolMemoPath = "tool"

func (c *Cache) toolDigest() []byte {
	tool.Lock()
	defer tool.Unlock()
	if tool.digest == nil {
		tool.digest = c.computeToolDigest()
	}
	return tool.digest
}

func (c *Cache) computeToolDigest() []byte {
	path, err := os.Executable()
	if err != nil {
		return []byte(info.Version)
	}
	fi, err := os.Stat(path)
	if err != nil {
		return []byte(info.Version)
	}

	var memo toolMemo
	if data, err := c.root.ReadFile(toolMemoPath); err == nil {
		if err := json.Unmarshal(data, &memo); err == nil &&
			memo.Path == path &&
			memo.SizeBytes == fi.Size() &&
			memo.MtimeNS == fi.ModTime().UnixNano() &&
			memo.SHA256 != "" {
			return []byte(memo.SHA256)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return []byte(info.Version)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return []byte(info.Version)
	}
	sum := hex.EncodeToString(h.Sum(nil))

	memo = toolMemo{
		Path:      path,
		SizeBytes: fi.Size(),
		MtimeNS:   fi.ModTime().UnixNano(),
		SHA256:    sum,
	}
	if data, err := json.Marshal(memo); err == nil {
		if f, name, err := c.CAS.createTemp("tool-"); err == nil {
			if _, werr := f.Write(data); werr == nil && f.Close() == nil {
				c.root.Rename(name, toolMemoPath)
			} else {
				f.Close()
			}
			c.root.Remove(name)
		}
	}
	return []byte(sum)
}
