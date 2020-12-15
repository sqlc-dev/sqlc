package golang

import (
	"fmt"
	"sort"
	"strings"
)

type Field struct {
	Name    string
	Type    string
	Tags    map[string]string
	Comment string
}

func (gf Field) Tag() string {
	tags := make([]string, 0, len(gf.Tags))
	for key, val := range gf.Tags {
		tags = append(tags, fmt.Sprintf("%s\"%s\"", key, val))
	}
	if len(tags) == 0 {
		return ""
	}
	sort.Strings(tags)
	return strings.Join(tags, " ")
}
