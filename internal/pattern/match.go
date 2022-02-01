package pattern

import (
	"fmt"
	"regexp"
)

// Match is a wrapper of *regexp.Regexp.
// It contains the match pattern compiled into a regular expression.
type Match struct {
	*regexp.Regexp
}

// Compile takes our match expression as a string, and compiles it into a *Match object.
// Will return an error on an invalid pattern.
func MatchCompile(pattern string) (match *Match, err error) {
	regex := ""
	escaped := false
	arr := []byte(pattern)

	for i := 0; i < len(arr); i++ {
		if escaped {
			escaped = false
			switch arr[i] {
			case '*', '?', '\\':
				regex += "\\" + string(arr[i])
			default:
				return nil, fmt.Errorf("Invalid escaped character '%c'", arr[i])
			}
		} else {
			switch arr[i] {
			case '\\':
				escaped = true
			case '*':
				regex += ".*"
			case '?':
				regex += "."
			case '.', '(', ')', '+', '|', '^', '$', '[', ']', '{', '}':
				regex += "\\" + string(arr[i])
			default:
				regex += string(arr[i])
			}
		}
	}

	if escaped {
		return nil, fmt.Errorf("Unterminated escape at end of pattern")
	}

	var r *regexp.Regexp

	if r, err = regexp.Compile("^" + regex + "$"); err != nil {
		return nil, err
	}

	return &Match{r}, nil
}
