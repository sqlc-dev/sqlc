package clickhouse

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TestArrayJoinColumnAliases validates that ARRAY JOIN creates properly aliased columns
// These columns should be available for reference in the SELECT list
func TestArrayJoinColumnAliases(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name             string
		query            string
		expectedColnames []string // column names from ARRAY JOIN
		wantErr          bool
	}{
		{
			name: "simple array join with alias",
			query: `
				SELECT id, tag
				FROM users
				ARRAY JOIN tags AS tag
			`,
			expectedColnames: []string{"tag"},
			wantErr:          false,
		},
		{
			name: "single array join with table alias and qualified name",
			query: `
				SELECT u.id, u.name, tag
				FROM users u
				ARRAY JOIN u.tags AS tag
			`,
			expectedColnames: []string{"tag"},
			wantErr:          false,
		},
		{
			name: "multiple array joins with aliases",
			query: `
				SELECT event_id, event_name, prop_key, prop_value
				FROM events
				ARRAY JOIN properties.keys AS prop_key, properties.values AS prop_value
			`,
			expectedColnames: []string{"prop_key", "prop_value"},
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := parser.Parse(strings.NewReader(tt.query))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse error: %v, wantErr %v", err, tt.wantErr)
			}

			if len(stmts) == 0 {
				t.Fatal("No statements parsed")
			}

			selectStmt, ok := stmts[0].Raw.Stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt, got %T", stmts[0].Raw.Stmt)
			}

			// Check that the FROM clause contains the ARRAY JOIN as a RangeSubselect
			if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) == 0 {
				t.Fatal("No FROM clause items found")
			}

			// Find the RangeSubselect that represents the ARRAY JOIN
			var arrayJoinRangeSubselect *ast.RangeSubselect
			for _, item := range selectStmt.FromClause.Items {
				if rs, ok := item.(*ast.RangeSubselect); ok {
					arrayJoinRangeSubselect = rs
					break
				}
			}

			if arrayJoinRangeSubselect == nil {
				t.Fatal("No RangeSubselect found for ARRAY JOIN")
			}

			// Verify that the RangeSubselect has a Subquery (synthetic SELECT statement)
			if arrayJoinRangeSubselect.Subquery == nil {
				t.Error("ARRAY JOIN RangeSubselect has no Subquery")
				return
			}

			syntheticSelect, ok := arrayJoinRangeSubselect.Subquery.(*ast.SelectStmt)
			if !ok {
				t.Errorf("Expected SelectStmt subquery, got %T", arrayJoinRangeSubselect.Subquery)
				return
			}

			// Verify the target list has the expected column names
			if syntheticSelect.TargetList == nil || len(syntheticSelect.TargetList.Items) == 0 {
				t.Error("Synthetic SELECT has no target list")
				return
			}

			if len(syntheticSelect.TargetList.Items) != len(tt.expectedColnames) {
				t.Errorf("Expected %d targets, got %d", len(tt.expectedColnames), len(syntheticSelect.TargetList.Items))
				return
			}

			// Verify the target values (which should be ResTargets with Name set)
			for i, expected := range tt.expectedColnames {
				target, ok := syntheticSelect.TargetList.Items[i].(*ast.ResTarget)
				if !ok {
					t.Errorf("Target %d is not a ResTarget: %T", i, syntheticSelect.TargetList.Items[i])
					continue
				}

				if target.Name == nil || *target.Name != expected {
					var name string
					if target.Name != nil {
						name = *target.Name
					}
					t.Errorf("Target %d: expected name %q, got %q", i, expected, name)
				}
			}
		})
	}
}
