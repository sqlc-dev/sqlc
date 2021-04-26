package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
)

type Field struct {
	Name    string
	Struct  string
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

func JSONTagName(name string, settings config.CombinedSettings) string {
	style := settings.Go.JSONTagsCaseStyle
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

func toSnakeCase(s string) string {
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
