package migrations

import (
	"bufio"
	"strings"
)

// Remove all lines after a rollback comment.
//
// goose:       -- +goose Down
// sql-migrate: -- +migrate Down
// tern:        ---- create above / drop below ----
// dbmate:      -- migrate:down
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
		if strings.HasPrefix(s.Text(), "---- create above / drop below ----") {
			break
		}
		if strings.HasPrefix(s.Text(), "-- migrate:down") {
			break
		}
		lines = append(lines, s.Text())
	}
	return strings.Join(lines, "\n")
}

func IsDown(filename string) bool {
	// Remove golang-migrate rollback files.
	return strings.HasSuffix(filename, ".down.sql")
}
