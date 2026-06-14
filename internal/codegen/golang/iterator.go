package golang

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/metadata"
)

const (
	IteratorScopeGlobal       = "global"
	IteratorScopeExplicitOnly = "explicit_only"

	IteratorStyleSeq2     = "seq2"
	IteratorStyleCallback = "callback"
	IteratorStyleRows     = "rows"

	IteratorStartLazy  = "lazy"
	IteratorStartEager = "eager"
)

func queryStreamAnnotated(comments []string) bool {
	for _, c := range comments {
		if c == metadata.StreamAnnotationComment || commentHasStreamAnnotation(c) {
			return true
		}
	}
	return false
}

func filterStreamAnnotationComments(comments []string) []string {
	if len(comments) == 0 {
		return comments
	}
	out := make([]string, 0, len(comments))
	for _, c := range comments {
		if c == metadata.StreamAnnotationComment {
			continue
		}
		out = append(out, c)
	}
	return out
}

func commentHasStreamAnnotation(s string) bool {
	return strings.Contains(s, metadata.CmdStream) ||
		strings.Contains(s, metadata.CmdManyStream)
}

func shouldEmitIterator(options *opts.Options, cmd string, comments []string) bool {
	if !options.EmitIterators || cmd != metadata.CmdMany {
		return false
	}
	switch options.IteratorScope {
	case "", IteratorScopeGlobal:
		return true
	case IteratorScopeExplicitOnly:
		return queryStreamAnnotated(comments)
	default:
		return false
	}
}

func iteratorMethodName(methodName, prefix string) string {
	if prefix == "" {
		prefix = "Iter"
	}
	for _, p := range []string{"List", "Get", "Find", "Select", "Fetch", "Stream"} {
		if strings.HasPrefix(methodName, p) {
			return prefix + strings.TrimPrefix(methodName, p)
		}
	}
	return prefix + methodName
}

func (v QueryValue) zeroValue() string {
	t := v.DefineType()
	if v.IsPointer() {
		return "nil"
	}
	return t + "{}"
}
