package shfmt

import (
	"os"
	"regexp"
	"strings"
)

var pat = regexp.MustCompile(`\$\{[A-Z_]+\}`)

type Replacer struct {
	envmap map[string]string
}

func (r *Replacer) Replace(f string) string {
	return pat.ReplaceAllStringFunc(f, func(s string) string {
		s = strings.TrimPrefix(s, "${")
		s = strings.TrimSuffix(s, "}")
		return r.envmap[s]
	})
}

func NewReplacer(env []string) *Replacer {
	r := Replacer{
		envmap: map[string]string{},
	}
	if env == nil {
		env = os.Environ()
	}
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		if k == "SQLC_AUTH_TOKEN" {
			continue
		}
		r.envmap[k] = v
	}
	return &r
}
