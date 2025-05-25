package selector

import "testing"

func TestDefaultSelectorColumnExpr(t *testing.T) {
	t.Parallel()

	selector := NewDefaultSelector()

	expectExpr := func(expected, name, dataType string) {
		if actual := selector.ColumnExpr(name, dataType); expected != actual {
			t.Errorf("Expected %v, got %v for data type %v", expected, actual, dataType)
		}
	}

	expectExpr("my_column", "my_column", "integer")
	expectExpr("my_column", "my_column", "json")
	expectExpr("my_column", "my_column", "jsonb")
	expectExpr("my_column", "my_column", "text")
}
