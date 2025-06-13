package compiler

import "testing"

func TestSelector(t *testing.T) {
	t.Parallel()

	selectorExpectColumnExpr := func(t *testing.T, selector selector, expected, name string, column *Column) {
		if actual := selector.ColumnExpr(name, column); expected != actual {
			t.Errorf("Expected %v, got %v for data type %v", expected, actual, column.DataType)
		}
	}

	t.Run("DefaultSelectorColumnExpr", func(t *testing.T) {
		t.Parallel()

		selector := newDefaultSelector()

		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "integer"})
		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "json"})
		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "jsonb"})
		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "text"})
	})

	t.Run("SQLiteSelectorColumnExpr", func(t *testing.T) {
		t.Parallel()

		selector := newSQLiteSelector()

		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "integer"})
		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "json"})
		selectorExpectColumnExpr(t, selector, "json(my_column)", "my_column", &Column{DataType: "jsonb"})
		selectorExpectColumnExpr(t, selector, "my_column", "my_column", &Column{DataType: "text"})
	})
}
