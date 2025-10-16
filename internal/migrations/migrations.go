package migrations

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
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
		statement := strings.ToLower(s.Text())
		if strings.HasPrefix(statement, "-- +goose down") {
			break
		}
		if strings.HasPrefix(statement, "-- +migrate down") {
			break
		}
		if strings.HasPrefix(statement, "---- create above / drop below ----") {
			break
		}
		if strings.HasPrefix(statement, "-- migrate:down") {
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

var ternTemplateRegex *regexp.Regexp

// tern: {{ template "filepath" . }}
func TransformStatements(pwd, content string) (string, error) {
	if !strings.Contains(content, "{{ template \"") {
		return content, nil
	}

	var err error
	var processed string
	
	if ternTemplateRegex == nil {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("failed to compile regexp: %v\n", r)
			}
			// It is tested, just recovering for it's technically possible to panic
		}()
		ternTemplateRegex = regexp.MustCompile(`\{\{ template \"(.+)\" \. \}\}`)
	}

	processed = ternTemplateRegex.ReplaceAllStringFunc(content, func(match string) string {
		filePath := ternTemplateRegex.FindStringSubmatch(match)[1]
		filePath = path.Join(pwd, filePath)
		read, err := os.ReadFile(filePath)
		if err != nil {
			err = errors.Join(err, fmt.Errorf("error reading file %s: %w", filePath, err))
			return match
		}
		return string(read)
	})

	return processed, err
}
