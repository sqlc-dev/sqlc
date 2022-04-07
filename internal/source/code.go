package source

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type Edit struct {
	Location int
	Old      string
	New      string
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
	for _, edit := range a {
		start := edit.Location
		if start > len(s) {
			return "", fmt.Errorf("edit start location is out of bounds")
		}
		if len(edit.New) <= 0 {
			return "", fmt.Errorf("empty edit contents")
		}
		if len(edit.Old) <= 0 {
			return "", fmt.Errorf("empty edit contents")
		}
		stop := edit.Location + len(edit.Old) - 1 // Assumes edit.New is non-empty
		if stop < len(s) {
			s = s[:start] + edit.New + s[stop+1:]
		} else {
			s = s[:start] + edit.New
		}
	}
	return s, nil
}

func StripComments(sql string) (string, []string, error) {
	re := regexp.MustCompile(`(?s)\/\*.*?\*\/\n?|--.*?\n`)
	re2 := regexp.MustCompile(`\/\*|\*\/|--|^\s*$`) // removes the comments
	// It will also remove the comments inside the comments
	sql = strings.TrimSpace(sql)

	commentsIndex := re.FindAllStringIndex(sql, -1)
	i := 0
	lines, comments := new(bytes.Buffer), new(bytes.Buffer)

	cleanComment := func(line string) string {
		if strings.HasPrefix(line, "-- name:") ||
			strings.HasPrefix(line, "/* name:") {
			return ""
		} else {
			return re2.ReplaceAllString(line, "")
		}
	}

	for _, comment := range commentsIndex {
		if i != comment[0] {
			lines.WriteString(sql[i:comment[0]])
		}
		comments.WriteString(cleanComment(sql[comment[0]:comment[1]]))
		i = comment[1]
	}

	if i < len(sql) {
		lines.WriteString(sql[i:])
	}
	commentsStr := comments.String()
	commentsStr = strings.TrimSuffix(commentsStr, "\n")
	return lines.String(), strings.SplitN(commentsStr, "\n", -1), nil
}
