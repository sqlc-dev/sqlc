package compiler

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
)

func TestSourceTableNames(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "cte, join and subquery dependencies",
			sql: `WITH filtered_accounts AS (
				SELECT account_id FROM accounts WHERE accounts.space_id = $1
				AND NOT EXISTS (
					SELECT 1 FROM account_tags t WHERE t.account_id = accounts.account_id
				)
			)
			SELECT acc.* FROM accounts acc
			JOIN filtered_accounts fa ON acc.account_id = fa.account_id
			LEFT JOIN transactions t ON t.debit_account_id = acc.account_id`,
			want: []string{"account_tags", "accounts", "transactions"},
		},
		{
			name: "deduplicated across aliases",
			sql:  `SELECT a1.account_id FROM accounts a1 JOIN accounts a2 ON a1.space_id = a2.space_id`,
			want: []string{"accounts"},
		},
		{
			name: "insert excludes write target, includes read",
			sql:  `INSERT INTO audit_log (account_id) SELECT account_id FROM accounts`,
			want: []string{"accounts"},
		},
		{
			name: "schema-qualified tables stay distinct",
			sql:  `SELECT 1 FROM audit.accounts JOIN accounts ON true`,
			want: []string{"accounts", "audit.accounts"},
		},
		{
			name: "no base tables",
			sql:  `SELECT 1`,
			want: []string{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := postgresql.NewParser().Parse(strings.NewReader(tc.sql + ";"))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(stmts) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(stmts))
			}
			got := sourceTableNames(stmts[0].Raw.Stmt)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("sourceTableNames\n  got:  %v\n  want: %v", got, tc.want)
			}
		})
	}
}
