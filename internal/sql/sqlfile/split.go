package sqlfile

import (
	"bufio"
	"context"
	"io"
	"strings"
)

// Split reads SQL queries from an io.Reader and returns them as a slice of strings.
// Each SQL query is delimited by a semicolon (;).
// The function handles:
// - Single-line comments (-- comment)
// - Multi-line comments (/* comment */)
// - Single-quoted strings ('string')
// - Double-quoted identifiers ("identifier")
// - Dollar-quoted strings ($$string$$ or $tag$string$tag$)
func Split(ctx context.Context, r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var queries []string
	var currentQuery strings.Builder
	var inSingleQuote bool
	var inDoubleQuote bool
	var inDollarQuote bool
	var dollarTag string
	var inMultiLineComment bool

	for scanner.Scan() {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		line := scanner.Text()
		i := 0
		lineLen := len(line)

		for i < lineLen {
			ch := line[i]

			// Handle multi-line comments
			if inMultiLineComment {
				if i+1 < lineLen && ch == '*' && line[i+1] == '/' {
					inMultiLineComment = false
					currentQuery.WriteString("*/")
					i += 2
					continue
				}
				currentQuery.WriteByte(ch)
				i++
				continue
			}

			// Handle dollar-quoted strings (PostgreSQL)
			if inDollarQuote {
				if ch == '$' {
					// Try to match the closing tag
					endTag := extractDollarTag(line[i:])
					if endTag == dollarTag {
						inDollarQuote = false
						currentQuery.WriteString(endTag)
						i += len(endTag)
						continue
					}
				}
				currentQuery.WriteByte(ch)
				i++
				continue
			}

			// Handle single-quoted strings
			if inSingleQuote {
				currentQuery.WriteByte(ch)
				if ch == '\'' {
					// Check for escaped quote ''
					if i+1 < lineLen && line[i+1] == '\'' {
						currentQuery.WriteByte('\'')
						i += 2
						continue
					}
					inSingleQuote = false
				}
				i++
				continue
			}

			// Handle double-quoted identifiers
			if inDoubleQuote {
				currentQuery.WriteByte(ch)
				if ch == '"' {
					// Check for escaped quote ""
					if i+1 < lineLen && line[i+1] == '"' {
						currentQuery.WriteByte('"')
						i += 2
						continue
					}
					inDoubleQuote = false
				}
				i++
				continue
			}

			// Check for single-line comment
			if i+1 < lineLen && ch == '-' && line[i+1] == '-' {
				// Rest of line is a comment
				currentQuery.WriteString(line[i:])
				break
			}

			// Check for multi-line comment start
			if i+1 < lineLen && ch == '/' && line[i+1] == '*' {
				inMultiLineComment = true
				currentQuery.WriteString("/*")
				i += 2
				continue
			}

			// Check for dollar quote start
			if ch == '$' {
				tag := extractDollarTag(line[i:])
				if tag != "" {
					inDollarQuote = true
					dollarTag = tag
					currentQuery.WriteString(tag)
					i += len(tag)
					continue
				}
			}

			// Check for single quote
			if ch == '\'' {
				inSingleQuote = true
				currentQuery.WriteByte(ch)
				i++
				continue
			}

			// Check for double quote
			if ch == '"' {
				inDoubleQuote = true
				currentQuery.WriteByte(ch)
				i++
				continue
			}

			// Check for semicolon (statement terminator)
			if ch == ';' {
				currentQuery.WriteByte(ch)
				// Check if there's a comment after the semicolon on the same line
				i++
				if i < lineLen {
					// Skip whitespace
					for i < lineLen && (line[i] == ' ' || line[i] == '\t') {
						currentQuery.WriteByte(line[i])
						i++
					}
					// If there's a comment, include it
					if i+1 < lineLen && line[i] == '-' && line[i+1] == '-' {
						currentQuery.WriteString(line[i:])
					}
				}
				query := strings.TrimSpace(currentQuery.String())
				if query != "" && query != ";" {
					queries = append(queries, query)
				}
				currentQuery.Reset()
				break // Move to next line
			}

			// Regular character
			currentQuery.WriteByte(ch)
			i++
		}

		// Add newline if we're building a query
		if currentQuery.Len() > 0 {
			currentQuery.WriteByte('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Handle any remaining query
	query := strings.TrimSpace(currentQuery.String())
	if query != "" && query != ";" {
		queries = append(queries, query)
	}

	return queries, nil
}

// extractDollarTag extracts a dollar-quoted string tag from the beginning of s.
// Returns empty string if no valid dollar tag is found.
// Valid tags: $$ or $identifier$ where identifier contains only alphanumeric and underscore.
func extractDollarTag(s string) string {
	if len(s) == 0 || s[0] != '$' {
		return ""
	}

	// Find the closing $
	for i := 1; i < len(s); i++ {
		if s[i] == '$' {
			tag := s[:i+1]
			// Validate tag content (only alphanumeric and underscore allowed between $)
			tagContent := tag[1 : len(tag)-1]
			if isValidDollarTagContent(tagContent) {
				return tag
			}
			return ""
		}
		// If we hit a character that's not allowed in a tag, it's not a dollar quote
		if !isValidDollarTagChar(s[i]) {
			return ""
		}
	}

	return ""
}

// isValidDollarTagContent returns true if s contains only valid characters for a dollar tag.
func isValidDollarTagContent(s string) bool {
	if s == "" {
		return true // $$ is valid
	}
	for _, ch := range s {
		if !isValidDollarTagChar(byte(ch)) {
			return false
		}
	}
	return true
}

// isValidDollarTagChar returns true if ch is a valid character in a dollar tag.
// Valid characters are alphanumeric and underscore.
func isValidDollarTagChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}
