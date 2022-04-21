package golang

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var IdentPattern = regexp.MustCompile("[^a-zA-Z0-9_]+")

type Constant struct {
	Name  string
	Type  string
	Value string
}

type Enum struct {
	Name      string
	Comment   string
	Constants []Constant
}

func EnumReplace(value string) string {
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	return IdentPattern.ReplaceAllString(id, "")
}

func EnumValueName(value string) string {
	name := ""
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	id = IdentPattern.ReplaceAllString(id, "")
	for _, part := range strings.Split(id, "_") {
		name += cases.Title(language.English).String(part)
	}
	return name
}
