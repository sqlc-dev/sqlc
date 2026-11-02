package lang

import "testing"

func TestIsJSONNullableOperator(t *testing.T) {
	t.Parallel()
	for _, op := range []string{"->", "->>", "#>", "#>>"} {
		if !IsJSONNullableOperator(op) {
			t.Errorf("expected %q to be classified as JSON-nullable", op)
		}
	}
	for _, op := range []string{"", "+", "-", "=", "::", "@>", "<@", "->>>"} {
		if IsJSONNullableOperator(op) {
			t.Errorf("did not expect %q to be classified as JSON-nullable", op)
		}
	}
}
