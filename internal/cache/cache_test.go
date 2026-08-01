package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"testing"
)

func testCache(t *testing.T) *Cache {
	t.Helper()
	c, err := OpenAt(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestCASRoundTrip(t *testing.T) {
	c := testCache(t)
	blob := []byte("SELECT id, name FROM authors")

	d, err := c.CAS.Put(blob)
	if err != nil {
		t.Fatal(err)
	}
	if d != DigestOf(blob) {
		t.Errorf("digest mismatch: %s != %s", d, DigestOf(blob))
	}
	if !c.CAS.Contains(d) {
		t.Error("Contains returned false for stored blob")
	}

	got, err := c.CAS.Get(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(blob) {
		t.Errorf("got %q, want %q", got, blob)
	}

	// Idempotent re-put
	if _, err := c.CAS.Put(blob); err != nil {
		t.Fatal(err)
	}
}

func TestCASMiss(t *testing.T) {
	c := testCache(t)
	if _, err := c.CAS.Get(DigestOf([]byte("never stored"))); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
	if c.CAS.Contains(Digest{Hash: "zz", SizeBytes: 1}) {
		t.Error("Contains returned true for invalid digest")
	}
}

func TestCASCorruptionEvicted(t *testing.T) {
	c := testCache(t)
	d, err := c.CAS.Put([]byte("original contents"))
	if err != nil {
		t.Fatal(err)
	}

	// Corrupt the blob on disk while keeping its size.
	path := c.CAS.path(d)
	if err := os.WriteFile(path, []byte("tampered contents"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := c.CAS.Get(d); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound for corrupt blob, got %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("corrupt blob was not evicted")
	}
}

func TestCASPutRepairsCorruptEntry(t *testing.T) {
	c := testCache(t)
	blob := []byte("plugin bytes")
	d, err := c.CAS.Put(blob)
	if err != nil {
		t.Fatal(err)
	}

	// Corrupt the entry on disk with different-sized contents, then Put the
	// correct blob again: the entry must be repaired, not trusted.
	path := c.CAS.path(d)
	if err := os.WriteFile(path, append(blob, "tampered"...), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CAS.Put(blob); err != nil {
		t.Fatal(err)
	}
	got, err := c.CAS.Get(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(blob) {
		t.Errorf("got %q, want %q", got, blob)
	}
}

func TestCASSHA256DirectLoad(t *testing.T) {
	c := testCache(t)
	blob := []byte("fetched plugin binary")
	s := sha256.Sum256(blob)
	declared := hex.EncodeToString(s[:])

	// A declared checksum whose content was never fetched is a miss.
	if _, err := c.CAS.Get(SHA256Digest(declared)); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound before Put, got %v", err)
	}

	d, err := c.CAS.PutSHA256(blob)
	if err != nil {
		t.Fatal(err)
	}
	if d.Hash != declared {
		t.Errorf("stored under %s, want declared checksum %s", d.Hash, declared)
	}

	// Later runs load the blob using only the checksum from the config file,
	// with no size and no action cache entry.
	got, err := c.CAS.Get(SHA256Digest(declared))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(blob) {
		t.Errorf("got %q, want %q", got, blob)
	}

	// sha256- and blake3-keyed entries live in separate namespaces.
	if c.CAS.Contains(DigestOf(blob)) {
		t.Error("blob stored by sha256 must not be visible under its blake3 digest")
	}
}

func TestCASSHA256CorruptionEvicted(t *testing.T) {
	c := testCache(t)
	blob := []byte("plugin contents")
	d, err := c.CAS.PutSHA256(blob)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(c.CAS.path(d), []byte("tampered!contents"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CAS.Get(SHA256Digest(d.Hash)); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound for corrupt blob, got %v", err)
	}
}

func TestActionDigestFraming(t *testing.T) {
	a := NewAction("Test").AddInput("query", []byte("ab")).AddInput("schema", []byte("c"))
	b := NewAction("Test").AddInput("query", []byte("a")).AddInput("schema", []byte("bc"))
	if a.Digest() == b.Digest() {
		t.Error("shifting bytes between inputs must change the action digest")
	}

	x := NewAction("Test").AddInput("query", []byte("q"))
	y := NewAction("Test").AddInput("query", []byte("q"))
	if x.Digest() != y.Digest() {
		t.Error("identical actions must have identical digests")
	}
}

func TestActionCacheRoundTrip(t *testing.T) {
	c := testCache(t)
	out, err := c.CAS.Put([]byte("analysis result"))
	if err != nil {
		t.Fatal(err)
	}

	action := NewAction("QueryAnalysis").AddInput("query", []byte("SELECT 1")).Digest()
	if _, err := c.Actions.Get(action); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound before Put, got %v", err)
	}

	if err := c.Actions.Put(action, &ActionResult{
		Outputs: map[string]Digest{"analysis.pb": out},
	}); err != nil {
		t.Fatal(err)
	}

	result, err := c.Actions.Get(action)
	if err != nil {
		t.Fatal(err)
	}
	if result.Outputs["analysis.pb"] != out {
		t.Errorf("got %s, want %s", result.Outputs["analysis.pb"], out)
	}
}

func TestActionCacheMissingOutputIsMiss(t *testing.T) {
	c := testCache(t)
	out, err := c.CAS.Put([]byte("ephemeral"))
	if err != nil {
		t.Fatal(err)
	}
	action := NewAction("Test").AddInput("in", []byte("x")).Digest()
	if err := c.Actions.Put(action, &ActionResult{
		Outputs: map[string]Digest{"out": out},
	}); err != nil {
		t.Fatal(err)
	}

	// Simulate the blob being garbage collected out from under the entry.
	if err := os.Remove(c.CAS.path(out)); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Actions.Get(action); !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound when output blob is gone, got %v", err)
	}
}

func TestActionCachePutRejectsMissingOutput(t *testing.T) {
	c := testCache(t)
	action := NewAction("Test").Digest()
	err := c.Actions.Put(action, &ActionResult{
		Outputs: map[string]Digest{"out": DigestOf([]byte("never stored"))},
	})
	if err == nil {
		t.Error("Put must reject results whose outputs are not in the CAS")
	}
}
