package dinosql

import (
	"bufio"
	"strings"
)

// Remove all lines after a rollback comment.
//
// goose:       -- +goose Down
// sql-migrate: -- +migrate Down
func RemoveRollbackStatements(contents string) string {
	s := bufio.NewScanner(strings.NewReader(contents))
	var lines []string
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "-- +goose Down") {
			break
		}
		if strings.HasPrefix(s.Text(), "-- +migrate Down") {
			break
		}
		lines = append(lines, s.Text())
	}
	return strings.Join(lines, "\n")
}
