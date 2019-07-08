package dinosql

import "strings"

// Remove all lines after a `-- +goose Down` comment
func RemoveGooseRollback(contents string) string {
	lines := strings.Split(contents, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "-- +goose Down") {
			lines = lines[:i]
			break
		}
	}
	return strings.Join(lines, "\n")
}
