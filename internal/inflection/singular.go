package inflection

import (
	"strings"

	upstream "github.com/jinzhu/inflection"
)

func Singular(name string) string {
	// Manual fix for incorrect handling of "campus"
	//
	// https://github.com/kyleconroy/sqlc/issues/430
	// https://github.com/jinzhu/inflection/issues/13
	if strings.ToLower(name) == "campus" {
		return name
	}
	return upstream.Singular(name)
}
