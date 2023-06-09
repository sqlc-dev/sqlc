package golang

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

func TagsToString(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}
	tagParts := make([]string, 0, len(tags))
	for key, val := range tags {
		tagParts = append(tagParts, fmt.Sprintf("%s:\"%s\"", key, val))
	}
	sort.Strings(tagParts)
	return strings.Join(tagParts, " ")
}

func JSONTagName(name string, settings *plugin.Settings) string {
	style := settings.Go.JsonTagsCaseStyle
	if style == "" || style == "none" {
		return name
	} else {
		return SetCaseStyle(name, style)
	}
}

func SetCaseStyle(name string, style string) string {
	switch style {
	case "camel":
		return toCamelCase(name)
	case "pascal":
		return toPascalCase(name)
	case "snake":
		return toSnakeCase(name)
	default:
		panic(fmt.Sprintf("unsupported JSON tags case style: '%s'", style))
	}
}

var camelPattern = regexp.MustCompile("[^A-Z][A-Z]+")

func toSnakeCase(s string) string {
	if !strings.ContainsRune(s, '_') {
		s = camelPattern.ReplaceAllStringFunc(s, func(x string) string {
			return x[:1] + "_" + x[1:]
		})
	}
	return strings.ToLower(s)
}

func toCamelCase(s string) string {
	return toCamelInitCase(s, false)
}

func toPascalCase(s string) string {
	return toCamelInitCase(s, true)
}

func toCamelInitCase(name string, initUpper bool) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if !initUpper && i == 0 {
			out += p
			continue
		}
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}
