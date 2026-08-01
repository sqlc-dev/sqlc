package cache

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"lukechampine.com/blake3"
)

// HashFunc identifies the hash function that keys a blob in the CAS.
//
// BLAKE3 is the native digest function. SHA256 exists so that content with an
// externally declared SHA-256 checksum — like a WASM plugin's checksum in
// sqlc's configuration — can be addressed directly by that checksum without
// any translation table.
type HashFunc string

const (
	BLAKE3 HashFunc = "blake3"
	SHA256 HashFunc = "sha256"
)

func sum(fn HashFunc, data []byte) (string, bool) {
	switch fn {
	case BLAKE3:
		s := blake3.Sum256(data)
		return hex.EncodeToString(s[:]), true
	case SHA256:
		s := sha256.Sum256(data)
		return hex.EncodeToString(s[:]), true
	}
	return "", false
}

// Digest identifies a blob by hash function, hash, and size, mirroring the
// Digest message from Bazel's remote execution API. The size is stored
// alongside the hash so that entries can be validated without reading blob
// contents; a negative size means the size is unknown, as with a checksum
// declared in a configuration file.
type Digest struct {
	// Function is the hash function; an empty value means BLAKE3.
	Function HashFunc `json:"function,omitempty"`
	// Hash is the lowercase hex-encoded 256-bit hash of the blob.
	Hash string `json:"hash"`
	// SizeBytes is the length of the blob in bytes, or negative if unknown.
	SizeBytes int64 `json:"size_bytes"`
}

func (d Digest) String() string {
	return fmt.Sprintf("%s:%s/%d", d.function(), d.Hash, d.SizeBytes)
}

func (d Digest) function() HashFunc {
	if d.Function == "" {
		return BLAKE3
	}
	return d.Function
}

func (d Digest) valid() bool {
	if _, ok := sum(d.function(), nil); !ok {
		return false
	}
	if len(d.Hash) != 64 {
		return false
	}
	_, err := hex.DecodeString(d.Hash)
	return err == nil
}

// DigestOf returns the BLAKE3 Digest of a blob.
func DigestOf(data []byte) Digest {
	hash, _ := sum(BLAKE3, data)
	return Digest{
		Hash:      hash,
		SizeBytes: int64(len(data)),
	}
}

// SHA256Digest returns a Digest referencing a blob by a declared SHA-256
// checksum whose size is not known, suitable for looking up remotely fetched
// content in the CAS.
func SHA256Digest(hexhash string) Digest {
	return Digest{
		Function:  SHA256,
		Hash:      hexhash,
		SizeBytes: -1,
	}
}

// An Action describes a unit of cacheable work, playing the role of Bazel's
// Action message: a mnemonic naming the kind of work plus the complete set of
// inputs that determine its outputs. Two actions with the same digest are
// assumed to produce the same outputs.
//
// Inputs are hashed incrementally with length-prefixed framing so that the
// boundary between inputs is unambiguous ("ab"+"c" hashes differently from
// "a"+"bc").
type Action struct {
	hasher *blake3.Hasher
}

// NewAction starts building an action key for the given mnemonic, e.g.
// "QueryAnalysis".
func NewAction(mnemonic string) *Action {
	a := &Action{hasher: blake3.New(32, nil)}
	a.write([]byte(mnemonic))
	return a
}

// AddInput mixes a named input into the action key. Order matters: callers
// must add inputs in a deterministic order.
func (a *Action) AddInput(name string, data []byte) *Action {
	a.write([]byte(name))
	a.write(data)
	return a
}

// Digest returns the action's digest, used as the action cache key.
func (a *Action) Digest() Digest {
	var sum [32]byte
	a.hasher.Sum(sum[:0])
	return Digest{
		Hash:      hex.EncodeToString(sum[:]),
		SizeBytes: 0,
	}
}

func (a *Action) write(data []byte) {
	var frame [8]byte
	binary.LittleEndian.PutUint64(frame[:], uint64(len(data)))
	a.hasher.Write(frame[:])
	a.hasher.Write(data)
}
