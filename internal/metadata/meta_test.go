package metadata

import "testing"

func TestNonMetadata(t *testing.T) {

	for _, query := range []string{
		`-- name: CreateFoo, :one`,
		`-- name: 9Foo_, :one`,
		`-- name: CreateFoo :two`,
		`-- name: CreateFoo`,
		`-- name: CreateFoo :one something`,
		`-- name: `,
		`--name: CreateFoo :one`,
		`--name CreateFoo :one`,
		`--name: CreateFoo :two`,
		"-- name:CreateFoo",
		`--name:CreateFoo :two`,
	} {
		if _, _, _, err := Parse(query, CommentSyntax{Dash: true}); err == nil {
			t.Errorf("expected invalid metadata: %q", query)
		}
	}

	for _, query := range []string{
		`-- some comment`,
		`-- name comment`,
		`--name comment`,
	} {
		if _, _, _, err := Parse(query, CommentSyntax{Dash: true}); err != nil {
			t.Errorf("expected valid comment: %q", query)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		query         string
		wantName      string
		wantType      string
		wantCmdParams CmdParams
	}{
		{
			query:    "-- name: CreateFoo :one",
			wantName: "CreateFoo",
			wantType: CmdOne,
		},
		{
			query:         "-- name: InsertMulti :exec multiple",
			wantName:      "InsertMulti",
			wantType:      CmdExec,
			wantCmdParams: CmdParams{InsertMultiple: true},
		},
		{
			query:         "-- name: SelectKey :many key=group_id",
			wantName:      "SelectKey",
			wantType:      CmdMany,
			wantCmdParams: CmdParams{ManyKey: "group_id"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			name, queryType, cmdParams, err := Parse(tc.query, CommentSyntax{Dash: true})
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if name != tc.wantName {
				t.Errorf("unexpected name: got %q; want %q", name, tc.wantName)
			}
			if queryType != tc.wantType {
				t.Errorf("unexpected queryType: got %q; want %q", queryType, tc.wantType)
			}
			if cmdParams != tc.wantCmdParams {
				t.Errorf("unexpected cmdParams: got %#v; want %#v", cmdParams, tc.wantCmdParams)
			}
		})
	}
}
