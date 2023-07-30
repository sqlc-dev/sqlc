package shfmt

import (
	"regexp"
	"strings"
)

var pat = regexp.MustCompile(`\$\{[A-Z_]+\}`)

func Replace(f string, vars map[string]string) string {
	return pat.ReplaceAllStringFunc(f, func(s string) string {
		s = strings.TrimPrefix(s, "${")
		s = strings.TrimSuffix(s, "}")
		return vars[s]
	})
}
