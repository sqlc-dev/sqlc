package golang

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/metadata"
)

func TestIteratorMethodName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		method string
		prefix string
		want   string
	}{
		{"ListAuthors", "Iter", "IterAuthors"},
		{"GetAuthors", "Iter", "IterAuthors"},
		{"FindAuthors", "Iter", "IterAuthors"},
		{"Authors", "Iter", "IterAuthors"},
		{"StreamAuthors", "Iter", "IterAuthors"},
	}
	for _, tc := range tests {
		got := iteratorMethodName(tc.method, tc.prefix)
		if got != tc.want {
			t.Errorf("iteratorMethodName(%q, %q) = %q, want %q", tc.method, tc.prefix, got, tc.want)
		}
	}
}

func TestQueryStreamAnnotated(t *testing.T) {
	t.Parallel()
	if !queryStreamAnnotated([]string{metadata.StreamAnnotationComment}) {
		t.Fatal("expected stream annotation comment")
	}
	if queryStreamAnnotated(nil) {
		t.Fatal("expected no annotation")
	}
}

func TestShouldEmitIterator(t *testing.T) {
	t.Parallel()
	global := &opts.Options{
		EmitIterators: true,
		IteratorScope: IteratorScopeGlobal,
	}
	if !shouldEmitIterator(global, metadata.CmdMany, nil) {
		t.Fatal("global scope should emit for :many")
	}
	if shouldEmitIterator(global, metadata.CmdOne, nil) {
		t.Fatal(":one should not emit iterator")
	}

	explicit := &opts.Options{
		EmitIterators: true,
		IteratorScope: IteratorScopeExplicitOnly,
	}
	if shouldEmitIterator(explicit, metadata.CmdMany, nil) {
		t.Fatal("explicit_only should skip unannotated :many")
	}
	if !shouldEmitIterator(explicit, metadata.CmdMany, []string{metadata.StreamAnnotationComment}) {
		t.Fatal("explicit_only should emit annotated stream query")
	}

	disabled := &opts.Options{EmitIterators: false, IteratorScope: IteratorScopeGlobal}
	if shouldEmitIterator(disabled, metadata.CmdMany, nil) {
		t.Fatal("emit_iterators false should not emit")
	}
}

func TestQueryValueZeroValue(t *testing.T) {
	t.Parallel()
	v := QueryValue{Struct: &Struct{Name: "Author"}}
	if v.zeroValue() != "Author{}" {
		t.Fatalf("got %q", v.zeroValue())
	}
	v.EmitPointer = true
	v.Struct = &Struct{Name: "Author"}
	if v.zeroValue() != "nil" {
		t.Fatalf("got %q", v.zeroValue())
	}
}
