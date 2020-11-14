package config

import (
	"fmt"
	"go/types"
	"strings"
)

type GoType struct {
	Path    string
	Package string
	Name    string
	Pointer bool
	BuiltIn bool
}

func ParseGoType(input string) (*GoType, error) {
	// validate GoType
	lastDot := strings.LastIndex(input, ".")
	lastSlash := strings.LastIndex(input, "/")
	typename := input
	var o GoType
	if lastDot == -1 && lastSlash == -1 {
		// if the type name has no slash and no dot, validate that the type is a basic Go type
		var found bool
		for _, typ := range types.Typ {
			info := typ.Info()
			if info == 0 {
				continue
			}
			if info&types.IsUntyped != 0 {
				continue
			}
			if typename == typ.Name() {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("Package override `go_type` specifier %q is not a Go basic type e.g. 'string'", input)
		}
		o.BuiltIn = true
	} else {
		// assume the type lives in a Go package
		if lastDot == -1 {
			return nil, fmt.Errorf("Package override `go_type` specifier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", input)
		}
		if lastSlash == -1 {
			return nil, fmt.Errorf("Package override `go_type` specifier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", input)
		}
		typename = input[lastSlash+1:]
		if strings.HasPrefix(typename, "go-") {
			// a package name beginning with "go-" will give syntax errors in
			// generated code. We should do the right thing and get the actual
			// import name, but in lieu of that, stripping the leading "go-" may get
			// us what we want.
			typename = typename[len("go-"):]
		}
		if strings.HasSuffix(typename, "-go") {
			typename = typename[:len(typename)-len("-go")]
		}
		o.Path = input[:lastDot]
	}
	o.Name = typename
	isPointer := input[0] == '*'
	if isPointer {
		o.Path = o.Path[1:]
		o.Name = "*" + o.Name
	}
	return &o, nil
}
