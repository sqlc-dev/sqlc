package config

import (
	"encoding/json"
	"fmt"
	"go/types"
	"regexp"
	"strings"
)

type GoType struct {
	Path    string `json:"import" yaml:"import"`
	Package string `json:"package" yaml:"package"`
	Name    string `json:"type" yaml:"type"`
	Pointer bool   `json:"pointer" yaml:"pointer"`
	Spec    string
	BuiltIn bool
}

type ParsedGoType struct {
	ImportPath string
	Package    string
	TypeName   string
	BasicType  bool
	StructTag  string
}

func (o *GoType) UnmarshalJSON(data []byte) error {
	var spec string
	if err := json.Unmarshal(data, &spec); err == nil {
		*o = GoType{Spec: spec}
		return nil
	}
	type alias GoType
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*o = GoType(a)
	return nil
}

func (o *GoType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var spec string
	if err := unmarshal(&spec); err == nil {
		*o = GoType{Spec: spec}
		return nil
	}
	type alias GoType
	var a alias
	if err := unmarshal(&a); err != nil {
		return err
	}
	*o = GoType(a)
	return nil
}

var validIdentifier = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var versionNumber = regexp.MustCompile(`^v[0-9]+$`)
var invalidIdentifier = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func generatePackageID(importPath string) (string, bool) {
	parts := strings.Split(importPath, "/")
	name := parts[len(parts)-1]
	// If the last part of the import path is a valid identifier, assume that's the package name
	if versionNumber.MatchString(name) && len(parts) >= 2 {
		name = parts[len(parts)-2]
		return invalidIdentifier.ReplaceAllString(strings.ToLower(name), "_"), true
	}
	if validIdentifier.MatchString(name) {
		return name, false
	}
	return invalidIdentifier.ReplaceAllString(strings.ToLower(name), "_"), true
}

// validate GoType
func (gt GoType) Parse() (*ParsedGoType, error) {
	var o ParsedGoType

	if gt.Spec == "" {
		// TODO: Validation
		if gt.Path == "" && gt.Package != "" {
			return nil, fmt.Errorf("Package override `go_type`: package name requires an import path")
		}
		var pkg string
		var pkgNeedsAlias bool

		if gt.Package == "" && gt.Path != "" {
			pkg, pkgNeedsAlias = generatePackageID(gt.Path)
			if pkgNeedsAlias {
				o.Package = pkg
			}
		} else {
			pkg = gt.Package
			o.Package = gt.Package
		}

		o.ImportPath = gt.Path
		o.TypeName = gt.Name
		o.BasicType = gt.Path == "" && gt.Package == ""
		if pkg != "" {
			o.TypeName = pkg + "." + o.TypeName
		}
		if gt.Pointer {
			o.TypeName = "*" + o.TypeName
		}
		return &o, nil
	}

	input := gt.Spec
	lastDot := strings.LastIndex(input, ".")
	lastSlash := strings.LastIndex(input, "/")
	typename := input
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
		o.BasicType = true
	} else {
		// assume the type lives in a Go package
		if lastDot == -1 {
			return nil, fmt.Errorf("Package override `go_type` specifier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", input)
		}
		typename = input[lastSlash+1:]
		// a package name beginning with "go-" will give syntax errors in
		// generated code. We should do the right thing and get the actual
		// import name, but in lieu of that, stripping the leading "go-" may get
		// us what we want.
		typename = strings.TrimPrefix(typename, "go-")
		typename = strings.TrimSuffix(typename, "-go")
		o.ImportPath = input[:lastDot]
	}
	o.TypeName = typename
	isPointer := input[0] == '*'
	if isPointer {
		o.ImportPath = o.ImportPath[1:]
		o.TypeName = "*" + o.TypeName
	}
	return &o, nil
}

// GoStructTag is a raw Go struct tag.
type GoStructTag string

// Parse parses and validates a GoStructTag.
// The output is in a form convenient for codegen.
//
// Sample valid inputs/outputs:
//
// In                    Out
// empty string      {}
// `a:"b"`           {"a": "b"}
// `a:"b" x:"y,z"`   {"a": "b", "x": "y,z"}
func (s GoStructTag) Parse() (map[string]string, error) {
	m := make(map[string]string)
	fields := strings.Fields(string(s))
	for _, f := range fields {
		k, v, ok := strings.Cut(f, ":")
		if !ok {
			return nil, fmt.Errorf("Failed to parse Go struct tag: no colon in field %q", f)
		}
		if len(v) < 2 || v[0] != '"' || v[len(v)-1] != '"' {
			return nil, fmt.Errorf("Failed to parse Go struct tag: missing quotes around value in field %q", f)
		}
		m[k] = v[1 : len(v)-1] // trim quotes off of v
	}
	return m, nil
}
