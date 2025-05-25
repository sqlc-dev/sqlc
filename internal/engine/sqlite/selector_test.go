package sqlite

import "testing"

func TestSelectorColumnExpr(t *testing.T) {
	t.Parallel()

	selector := NewSelector()

	expectExpr := func(expected, name, dataType string) {
		if actual := selector.ColumnExpr(name, dataType); expected != actual {
			t.Errorf("Expected %v, got %v for data type %v", expected, actual, dataType)
		}
	}

	expectExpr("my_column", "my_column", "integer")
	expectExpr("my_column", "my_column", "json")
	expectExpr("json(my_column)", "my_column", "jsonb")
	expectExpr("my_column", "my_column", "text")
}
