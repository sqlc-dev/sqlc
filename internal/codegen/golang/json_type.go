package golang

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// JSONType is a plain Go struct synthesized from a JSON directive. No
// Scan/Value methods are needed: pgx v5 decodes the jsonb value directly into
// a bare JSONType (or []JSONType for an array) via reflection.
type JSONType struct {
	Name   string
	Fields []Field
}

// newGoJSONColumn resolves the Go type for a JSON column ([]Name for arrays,
// Name otherwise) and collects the JSONType declarations it needs, including
// nested ones. col.JsonName is used as-is so two queries can share a type.
func newGoJSONColumn(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column, models modelTypeSet, qualifier string) (string, []JSONType) {
	structName := StructName(col.JsonName, options)

	var fields []Field
	var types []JSONType
	for _, f := range col.JsonFields {
		tags := map[string]string{"json": f.Name}
		addExtraGoStructTags(tags, req, options, f)

		var fieldType string
		if len(f.JsonFields) > 0 {
			var nested []JSONType
			fieldType, nested = newGoJSONColumn(req, options, f, models, qualifier)
			types = append(types, nested...)
		} else {
			fieldType = qualifyType(goType(req, options, f), models, qualifier)
		}

		// A namespaced rename ("Name.fieldKey") wins over the global one.
		fieldName, ok := options.Rename[col.JsonName+"."+f.Name]
		if !ok {
			fieldName = StructName(f.Name, options)
		}

		fields = append(fields, Field{
			Name:   fieldName,
			DBName: f.Name,
			Type:   fieldType,
			Tags:   tags,
			Column: f,
		})
	}

	types = append(types, JSONType{Name: structName, Fields: fields})

	if col.IsArray {
		return "[]" + structName, types
	}
	return structName, types
}

// internJSONTypes deduplicates types against seen (package-wide, keyed by
// name, mutated in place), returning only those still needing to be emitted.
// A name that collides with a model/enum, or with an earlier type of a
// different shape, is an error. Scalar and array uses share one declaration.
func internJSONTypes(seen map[string]JSONType, models modelTypeSet, types []JSONType) ([]JSONType, error) {
	var toEmit []JSONType
	for _, t := range types {
		if _, ok := models[t.Name]; ok {
			return nil, fmt.Errorf("sqlc.jsonb_build_object.%q(...) collides with an existing model or enum type; give it a different name", t.Name)
		}
		prev, ok := seen[t.Name]
		if !ok {
			seen[t.Name] = t
			toEmit = append(toEmit, t)
			continue
		}
		if !sameJSONShape(prev, t) {
			return nil, fmt.Errorf("sqlc.jsonb_build_object.%q(...) is used with conflicting shapes across queries; give one of them a different name", t.Name)
		}
	}
	return toEmit, nil
}

func sameJSONShape(a, b JSONType) bool {
	if len(a.Fields) != len(b.Fields) {
		return false
	}
	for i, f := range a.Fields {
		g := b.Fields[i]
		if f.Name != g.Name || f.Type != g.Type || f.Tag() != g.Tag() {
			return false
		}
	}
	return true
}
