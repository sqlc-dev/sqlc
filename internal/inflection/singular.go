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
	// https://github.com/sqlc-dev/sqlc/issues/430
	// https://github.com/jinzhu/inflection/issues/13
	if strings.ToLower(s.Name) == "campus" {
		return s.Name
	}
	// Manual fix for incorrect handling of "meta"
	//
	// https://github.com/sqlc-dev/sqlc/issues/1217
	// https://github.com/jinzhu/inflection/issues/21
	if strings.ToLower(s.Name) == "meta" {
		return s.Name
	}
	// Manual fix for incorrect handling of "calories"
	//
	// https://github.com/sqlc-dev/sqlc/issues/2017
	// https://github.com/jinzhu/inflection/issues/23
	if strings.ToLower(s.Name) == "calories" {
		return "calorie"
	}
	// Manual fix for incorrect handling of "-ves" suffix
	if strings.ToLower(s.Name) == "waves" {
		return "wave"
	}

	if strings.ToLower(s.Name) == "metadata" {
		return "metadata"
	}

	return upstream.Singular(s.Name)
}
