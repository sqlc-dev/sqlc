package dinosql

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	pg "github.com/lfittl/pg_query_go"
)

func TestValidateParamRef(t *testing.T) {
	// equateErrorMessage reports errors to be equal if both are nil
	// or both have the same message.
	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	for _, tc := range []struct {
		query string
		err   error
	}{
		{
			"SELECT foo FROM bar WHERE baz = $4;",
			fmt.Errorf("missing parameter reference: $1"),
		},
		{
			"SELECT foo FROM bar WHERE baz = $1;",
			nil,
		},
	} {
		tree, err := pg.Parse(tc.query)
		if err != nil {
			t.Fatal(err)
		}
		actual := validateParamRef(tree.Statements[0])
		if diff := cmp.Diff(tc.err, actual, equateErrorMessage); diff != "" {
			t.Errorf("error mismatch: \n%s", diff)
		}
	}
}
