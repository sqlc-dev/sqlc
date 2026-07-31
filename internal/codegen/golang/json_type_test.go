package golang

import (
	"strings"
	"testing"
)

func field(name, typ string) Field {
	return Field{Name: name, Type: typ, Tags: map[string]string{"json": name}}
}

func TestInternJSONTypes(t *testing.T) {
	t.Run("first occurrence is emitted", func(t *testing.T) {
		seen := map[string]JSONType{}
		toEmit, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{
			{Name: "Book", Fields: []Field{field("id", "int64")}},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(toEmit) != 1 {
			t.Fatalf("expected 1 type to emit, got %d", len(toEmit))
		}
	})

	t.Run("same name and shape is interned, not re-emitted", func(t *testing.T) {
		seen := map[string]JSONType{}
		book := JSONType{Name: "Book", Fields: []Field{field("id", "int64")}}
		if _, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{book}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		toEmit, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{book})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(toEmit) != 0 {
			t.Errorf("expected 0 types to emit on reuse, got %d", len(toEmit))
		}
	})

	t.Run("same name, different fields is a conflict", func(t *testing.T) {
		seen := map[string]JSONType{}
		if _, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{
			{Name: "Book", Fields: []Field{field("id", "int64")}},
		}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{
			{Name: "Book", Fields: []Field{field("id", "int64"), field("title", "string")}},
		})
		if err == nil || !strings.Contains(err.Error(), "conflicting shapes") {
			t.Errorf("error = %v, want mention of conflicting shapes", err)
		}
	})

	t.Run("scalar and array uses of the same name intern together", func(t *testing.T) {
		// The struct declaration is identical either way; only the
		// referencing field gains a [] prefix, so these aren't a conflict.
		seen := map[string]JSONType{}
		book := JSONType{Name: "Book", Fields: []Field{field("id", "int64")}}
		if _, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{book}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := internJSONTypes(seen, modelTypeSet{}, []JSONType{book}); err != nil {
			t.Errorf("unexpected error reusing the same shape: %v", err)
		}
	})

	t.Run("name collides with an existing model", func(t *testing.T) {
		seen := map[string]JSONType{}
		models := modelTypeSet{"Book": struct{}{}}
		_, err := internJSONTypes(seen, models, []JSONType{
			{Name: "Book", Fields: []Field{field("id", "int64")}},
		})
		if err == nil || !strings.Contains(err.Error(), "collides with an existing model") {
			t.Errorf("error = %v, want mention of model collision", err)
		}
	})
}

func TestSameJSONShape(t *testing.T) {
	tests := []struct {
		name string
		a, b JSONType
		want bool
	}{
		{
			name: "identical",
			a:    JSONType{Fields: []Field{field("x", "int32")}},
			b:    JSONType{Fields: []Field{field("x", "int32")}},
			want: true,
		},
		{
			name: "different field count",
			a:    JSONType{Fields: []Field{field("x", "int32")}},
			b:    JSONType{Fields: []Field{field("x", "int32"), field("y", "string")}},
			want: false,
		},
		{
			name: "different field type",
			a:    JSONType{Fields: []Field{field("x", "int32")}},
			b:    JSONType{Fields: []Field{field("x", "string")}},
			want: false,
		},
		{
			name: "different field name",
			a:    JSONType{Fields: []Field{field("x", "int32")}},
			b:    JSONType{Fields: []Field{field("y", "int32")}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameJSONShape(tt.a, tt.b); got != tt.want {
				t.Errorf("sameJSONShape() = %v, want %v", got, tt.want)
			}
		})
	}
}
