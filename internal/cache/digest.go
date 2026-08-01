package cache

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"lukechampine.com/blake3"
)

// Digest identifies a blob by its BLAKE3 hash and size, mirroring the Digest
// message from Bazel's remote execution API. The size is stored alongside the
// hash so that entries can be validated without reading blob contents.
type Digest struct {
	// Hash is the lowercase hex-encoded BLAKE3-256 hash of the blob.
	Hash string `json:"hash"`
	// SizeBytes is the length of the blob in bytes.
	SizeBytes int64 `json:"size_bytes"`
}

func (d Digest) String() string {
	return fmt.Sprintf("blake3:%s/%d", d.Hash, d.SizeBytes)
}

func (d Digest) valid() bool {
	if len(d.Hash) != 64 || d.SizeBytes < 0 {
		return false
	}
	_, err := hex.DecodeString(d.Hash)
	return err == nil
}

// DigestOf returns the Digest of a blob.
func DigestOf(data []byte) Digest {
	sum := blake3.Sum256(data)
	return Digest{
		Hash:      hex.EncodeToString(sum[:]),
		SizeBytes: int64(len(data)),
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
// "QueryAnalysis" or "FetchPlugin".
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
