package mysql

import (
	"fmt"
	"regexp"
	"strconv"
)

func locFromSyntaxErr(errMessage error) (int, error) {
	matcher := regexp.MustCompile("position ([0-9]*)")
	results := matcher.FindStringSubmatch(errMessage.Error())
	if len(results) > 0 {
		return strconv.Atoi(results[1])
	}
	return 0, fmt.Errorf("failed to find position integer in parser error message")
}

func nearStrFromSyntaxErr(errMessage error) (string, error) {
	matcher := regexp.MustCompile("near '(.*)'")
	results := matcher.FindStringSubmatch(errMessage.Error())
	if len(results) > 0 {
		return results[1], nil
	}
	return "", fmt.Errorf("failed to find parser 'near' message")
}

type PositionedErr struct {
	Pos int
	Err error
}

func (e PositionedErr) Error() string {
	return e.Err.Error()
}
