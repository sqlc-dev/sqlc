package inflection

import (
	"strings"

	upstream "github.com/jinzhu/inflection"
)

type SingularParams struct {
	Name       string
	Exclusions []string
}

func Singular(s SingularParams) string {
	for _, exclusion := range s.Exclusions {
		if strings.EqualFold(s.Name, exclusion) {
			return s.Name
		}
	}

	// Manual fix for incorrect handling of "campus"
	//
	// https://github.com/kyleconroy/sqlc/issues/430
	// https://github.com/jinzhu/inflection/issues/13
	if strings.ToLower(s.Name) == "campus" {
		return s.Name
	}
	// Manual fix for incorrect handling of "meta"
	//
	// https://github.com/kyleconroy/sqlc/issues/1217
	// https://github.com/jinzhu/inflection/issues/21
	if strings.ToLower(s.Name) == "meta" {
		return s.Name
	}
	return upstream.Singular(s.Name)
}
