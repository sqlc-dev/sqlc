package golang

import (
	"fmt"
	"strings"
)

func ApplySchema(query string) string {
	tables := make(map[string]bool)
	ctes := make(map[string]bool)

	words := strings.Fields(query)

	// Getting all the table names and CTEs
	withinCTE := false
	for i, word := range words {
		upperWord := strings.ToUpper(word)

		if upperWord == "WITH" {
			withinCTE = true
			continue
		} else if withinCTE {
			ctes[words[i]] = true
			withinCTE = false
			continue
		}

		if isSQLKeyword(upperWord) {
			tables[nextNonKeyword(words, i)] = true
		}
	}

	// Removing from tables the CTEs
	for cte := range ctes {
		delete(tables, cte)
	}

	// Replacing the table names with the placeholder
	for table := range tables {
		query = strings.ReplaceAll(query, " "+table, fmt.Sprintf(" `%%s`.%s", table))
	}

	return query
}

// Helper function to check if a word is a relevant SQL keyword
func isSQLKeyword(word string) bool {
	switch word {
	case "FROM", "JOIN", "LEFT JOIN", "RIGHT JOIN", "FULL JOIN", "INNER JOIN", "CROSS JOIN", "UPDATE", "DELETE FROM", "INSERT INTO":
		return true
	}
	return false
}

func nextNonKeyword(words []string, currentIndex int) string {
	for i := currentIndex + 1; i < len(words); i++ {
		if !isSQLKeyword(words[i]) && words[i] != "AS" && words[i] != "(" {
			return words[i]
		}
	}
	return ""
}
