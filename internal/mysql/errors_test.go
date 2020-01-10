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
	}

	location, err := locFromSyntaxErr(parseErr)
	if err != nil {
		t.Errorf("failed to parse location from sqlparser syntax error message: %v", err)
	} else if location != expectedLocation {
		t.Errorf("parsed incorrect location from sqlparser syntax error message: %v", cmp.Diff(expectedLocation, location))
	}

	near, err := nearStrFromSyntaxErr(parseErr)
	if err != nil {
		t.Errorf("failed to parse 'nearby' chars from sqlparser syntax error message: %v", err)
	} else if near != expectedNear {
		t.Errorf("parse incorrect 'nearby' chars from sqlparser syntax error message: %v", cmp.Diff(expectedNear, near))
	}
}
