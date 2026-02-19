package pattern

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Match is a wrapper of *regexp.Regexp.
// It contains the match pattern compiled into a regular expression.
type Match struct {
	*regexp.Regexp
}

var (
	matchCache     = make(map[string]*Match)
	matchCacheLock sync.RWMutex
)

// Compile takes our match expression as a string, and compiles it into a *Match object.
// Will return an error on an invalid pattern.
func MatchCompile(pattern string) (*Match, error) {
	// check for pattern in cache
	matchCacheLock.RLock()
	matcher, ok := matchCache[pattern]
	matchCacheLock.RUnlock()
	if ok {
		return matcher, nil
	}

	// pattern isn't in cache, compile it
	matcher, err := matchCompile(pattern)
	if err != nil {
		return nil, err
	}
	// add it to the cache
	matchCacheLock.Lock()
	matchCache[pattern] = matcher
	matchCacheLock.Unlock()

	return matcher, nil
}

func matchCompile(pattern string) (match *Match, err error) {
	var regex strings.Builder
	escaped := false
	arr := []byte(pattern)

	for i := range arr {
		if escaped {
			escaped = false
			switch arr[i] {
			case '*', '?', '\\':
				regex.WriteString("\\" + string(arr[i]))
			default:
				return nil, fmt.Errorf("Invalid escaped character '%c'", arr[i])
			}
		} else {
			switch arr[i] {
			case '\\':
				escaped = true
			case '*':
				regex.WriteString(".*")
			case '?':
				regex.WriteString(".")
			case '.', '(', ')', '+', '|', '^', '$', '[', ']', '{', '}':
				regex.WriteString("\\" + string(arr[i]))
			default:
				regex.WriteString(string(arr[i]))
			}
		}
	}

	if escaped {
		return nil, fmt.Errorf("Unterminated escape at end of pattern")
	}

	var r *regexp.Regexp

	if r, err = regexp.Compile("^" + regex.String() + "$"); err != nil {
		return nil, err
	}

	return &Match{r}, nil
}
