package mysql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestArgName(t *testing.T) {
	tcase := [...]struct {
		input  string
		output string
	}{
		{
			input:  "get_users",
			output: "getUsers",
		},
		{
			input:  "get_users_by_id",
			output: "getUsersByID",
		},
		{
			input:  "get_all_",
			output: "getAll",
		},
	}

	for _, tc := range tcase {
		name := argName(tc.input)
		if diff := cmp.Diff(name, tc.output); diff != "" {
			t.Errorf(diff)
		}
	}
}
