package source

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type Edit struct {
	Location int
	Old      string
	New      string
	OldFunc  func(string) int
}

type CommentSyntax struct {
	Dash      bool
	Hash      bool
	SlashStar bool
}

func LineNumber(source string, head int) (int, int) {
	// Calculate the true line and column number for a query, ignoring spaces
	var comment bool
	var loc, line, col int
	for i, char := range source {
		loc += 1
		col += 1
		// TODO: Check bounds
		if char == '-' && source[i+1] == '-' {
			comment = true
		}
		if char == '\n' {
			comment = false
			line += 1
			col = 0
		}
		if loc <= head {
			continue
		}
		if unicode.IsSpace(char) {
			continue
		}
		if comment {
			continue
		}
		break
	}
	return line + 1, col
}

func Pluck(source string, location, length int) (string, error) {
	head := location
	tail := location + length
	return source[head:tail], nil
}

func Mutate(raw string, a []Edit) (string, error) {
	if len(a) == 0 {
		return raw, nil
	}

	sort.Slice(a, func(i, j int) bool { return a[i].Location > a[j].Location })

	s := raw
	for idx, edit := range a {
		start := edit.Location
		if start > len(s) || start < 0 {
			return "", fmt.Errorf("edit start location is out of bounds")
		}
		var oldLen int
		if edit.OldFunc != nil {
			oldLen = edit.OldFunc(s[start:])
		} else {
			oldLen = len(edit.Old)
		}

		stop := edit.Location + oldLen
		if stop > len(s) {
			return "", fmt.Errorf("edit stop location is out of bounds")
		}

		// If this is not the first edit, (applied backwards), check if
		// this edit overlaps the previous one (and is therefore a developer error)
		if idx != 0 {
			prevEdit := a[idx-1]
			if prevEdit.Location < edit.Location+oldLen {
				return "", fmt.Errorf("2 edits overlap")
			}
		}

		s = s[:start] + edit.New + s[stop:]
	}
	return s, nil
}

func StripComments(sql string) (string, []string, error) {
	s := bufio.NewScanner(strings.NewReader(strings.TrimSpace(sql)))
	var lines, comments []string
	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "-- name:") {
			continue
		}
		if strings.HasPrefix(t, "/* name:") && strings.HasSuffix(t, "*/") {
			continue
		}
		if strings.HasPrefix(t, "# name:") {
			continue
		}
		if strings.HasPrefix(t, "--") {
			comments = append(comments, strings.TrimPrefix(t, "--"))
			continue
		}
		if strings.HasPrefix(t, "/*") && strings.HasSuffix(t, "*/") {
			t = strings.TrimPrefix(t, "/*")
			t = strings.TrimSuffix(t, "*/")
			comments = append(comments, t)
			continue
		}
		if strings.HasPrefix(t, "#") {
			comments = append(comments, strings.TrimPrefix(t, "#"))
			continue
		}
		lines = append(lines, t)
	}
	return strings.Join(lines, "\n"), comments, s.Err()
}

func CleanedComments(rawSQL string, cs CommentSyntax) ([]string, error) {
	s := bufio.NewScanner(strings.NewReader(strings.TrimSpace(rawSQL)))
	var comments []string
	for s.Scan() {
		line := s.Text()
		var prefix string
		if strings.HasPrefix(line, "--") {
			if !cs.Dash {
				continue
			}
			prefix = "--"
		}
		if strings.HasPrefix(line, "/*") {
			if !cs.SlashStar {
				continue
			}
			prefix = "/*"
		}
		if strings.HasPrefix(line, "#") {
			if !cs.Hash {
				continue
			}
			prefix = "#"
		}
		if prefix == "" {
			continue
		}

		rest := line[len(prefix):]
		rest = strings.TrimSuffix(rest, "*/")
		comments = append(comments, rest)
	}
	return comments, s.Err()
}
