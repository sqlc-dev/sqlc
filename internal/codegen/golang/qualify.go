package golang

import "strings"

// stripQualifier removes a leading slice/pointer prefix and the given
// `pkg.` qualifier from a Go type expression. When the qualifier is empty
// or absent from the type, the input is returned unchanged.
func stripQualifier(t, qualifier string) string {
	if qualifier == "" {
		return t
	}
	prefix := ""
	rest := t
	for {
		if strings.HasPrefix(rest, "[]") {
			prefix += "[]"
			rest = rest[2:]
			continue
		}
		if strings.HasPrefix(rest, "*") {
			prefix += "*"
			rest = rest[1:]
			continue
		}
		break
	}
	if strings.HasPrefix(rest, qualifier) {
		return prefix + rest[len(qualifier):]
	}
	return t
}

// modelTypeSet is the set of Go type names that live in the models file.
type modelTypeSet map[string]struct{}

// buildModelTypeSet returns the set of type names that are declared in
// models.go for the current codegen invocation.
func buildModelTypeSet(enums []Enum, structs []Struct) modelTypeSet {
	set := make(modelTypeSet, len(enums)*4+len(structs))
	for _, e := range enums {
		set[e.Name] = struct{}{}
		set["Null"+e.Name] = struct{}{}
	}
	for _, s := range structs {
		if s.IsModel {
			set[s.Name] = struct{}{}
		}
	}
	return set
}

// qualifyType prefixes a Go type expression with `qualifier` when the bare
// type name belongs to `models`. Slice and pointer prefixes are preserved.
// When qualifier is empty (i.e. models live in the queries package), the
// input is returned unchanged.
func qualifyType(t string, models modelTypeSet, qualifier string) string {
	if qualifier == "" || t == "" || len(models) == 0 {
		return t
	}
	prefix := ""
	rest := t
	for {
		if strings.HasPrefix(rest, "[]") {
			prefix += "[]"
			rest = rest[2:]
			continue
		}
		if strings.HasPrefix(rest, "*") {
			prefix += "*"
			rest = rest[1:]
			continue
		}
		break
	}
	if _, ok := models[rest]; ok {
		return prefix + qualifier + rest
	}
	return t
}
