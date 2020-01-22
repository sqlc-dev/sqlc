package mysql

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"vitess.io/vitess/go/vt/sqlparser"
)

func TestSyntaxErr(t *testing.T) {
	tokenizer := sqlparser.NewStringTokenizer("SELEC T id FROM users;")
	expectedLocation := 6
	expectedNear := "SELEC"

	_, parseErr := sqlparser.ParseNextStrictDDL(tokenizer)
	if parseErr == nil {
		t.Errorf("Tokenizer failed to error on invalid MySQL syntax")
	} else if posErr, ok := parseErr.(sqlparser.PositionedErr); ok {
		if posErr.Pos != expectedLocation {
			t.Errorf(cmp.Diff(posErr.Pos, expectedLocation))
		}
		if string(posErr.Near) != expectedNear {
			t.Errorf(cmp.Diff(string(posErr.Near), string(expectedNear)))
		}
	} else {
		t.Errorf("failed to return sqlparser.PositionedErr error for invalid mysql expression")
	}
}
