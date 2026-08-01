package cache

import (
	"crypto/sha256"
	"io"
	"os"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/info"
)

// toolDigest returns the sha256 of the running sqlc binary, hashed once per
// process. The binary is an input to every action — a rebuilt sqlc may
// analyze queries or embed a different wazero than the one that produced a
// cache entry, even when the version string is unchanged (dev builds). If
// the executable can't be read, fall back to the version string.
var toolDigest = sync.OnceValue(func() []byte {
	path, err := os.Executable()
	if err != nil {
		return []byte(info.Version)
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
	return h.Sum(nil)
})
