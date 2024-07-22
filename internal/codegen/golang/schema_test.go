package golang_test

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang"
)

func TestApplySchema(t *testing.T) {
	testCases := []struct {
		name          string
		inputQuery    string
		expectedQuery string
	}{
		{
			name:          "Simple Query with Single Table",
			inputQuery:    "SELECT * FROM users",
			expectedQuery: "SELECT * FROM `%s`.users",
		},
		{
			name:          "Query with Multiple Tables",
			inputQuery:    "SELECT * FROM users JOIN orders ON users.id = orders.user_id",
			expectedQuery: "SELECT * FROM `%s`.users JOIN `%s`.orders ON `%s`.users.id = `%s`.orders.user_id",
		},
		{
			name:          "Query with CTE",
			inputQuery:    "WITH user_orders AS (SELECT * FROM users JOIN orders ON users.id = orders.user_id) SELECT * FROM user_orders",
			expectedQuery: "WITH user_orders AS (SELECT * FROM `%s`.users JOIN `%s`.orders ON `%s`.users.id = `%s`.orders.user_id) SELECT * FROM user_orders",
		},
		{
			name:          "Query with CTE and Aliases",
			inputQuery:    "WITH user_orders AS (SELECT * FROM users u JOIN orders o ON u.id = o.user_id) SELECT * FROM user_orders uo",
			expectedQuery: "WITH user_orders AS (SELECT * FROM `%s`.users u JOIN `%s`.orders o ON u.id = o.user_id) SELECT * FROM user_orders uo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := golang.ApplySchema(tc.inputQuery)
			if result != tc.expectedQuery {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expectedQuery, result)
			}
		})
	}
}

// "SELECT * FROM `%s`.users JOIN `%s`.orders ON `%s`.users.id `%s`.= `%s`.orders.user_id"
