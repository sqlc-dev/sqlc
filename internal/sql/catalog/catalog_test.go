package catalog

import "testing"

func TestNew(t *testing.T) {

	const defaultSchema = "default"

	newCatalog := New(defaultSchema)

	if newCatalog.DefaultSchema == "" {
		t.Errorf("newCatalog.DefaultSchema: want %s, got %s", defaultSchema, newCatalog.DefaultSchema)
	}

	if newCatalog.Schemas == nil {
		t.Error("newCatalog.Schemas should not be nil")
	}

	if len(newCatalog.Schemas) != 1 {
		t.Errorf("newCatalog.Schemas length want 1, got %d", len(newCatalog.Schemas))
	}

	if newCatalog.Schemas[0].Name != defaultSchema {
		t.Error("newCatalog.Schemas should have the default schema")
	}

	if newCatalog.Extensions == nil {
		t.Error("newCatalog.Extensions should be empty")
	}
}
