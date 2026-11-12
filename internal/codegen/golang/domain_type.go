package golang

import (
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// domainColumn is col with every domain its type names replaced by the
// domain's base, or nil when its type names none.
func domainColumn(req *plugin.GenerateRequest, col *plugin.Column) *plugin.Column {
	if col.TypeExpr != nil {
		t, ok := resolveDomains(req, col.TypeExpr)
		if !ok {
			return nil
		}
		out := proto.Clone(col).(*plugin.Column)
		out.TypeExpr = t
		name, dims := flatType(t)
		out.Type = &plugin.Identifier{Name: name}
		out.IsArray, out.ArrayDims = dims > 0, dims
		return out
	}
	d := findDomain(req, col.Type.GetSchema(), col.Type.GetName())
	if baseless(req, d) {
		return nil
	}
	out := proto.Clone(col).(*plugin.Column)
	out.Type = &plugin.Identifier{Name: d.Base.Name}
	return out
}

// flatType is t as the flat description reads it, the way
// coreResultCatalog writes one: the innermost element's name, with one
// dimension per array around it.
func flatType(t *plugin.TypeExpr) (string, int32) {
	var dims int32
	for typeExprName(t) == "array" && typeExprArg(t, 0) != nil {
		t = typeExprArg(t, 0)
		dims++
	}
	return t.GetName(), dims
}

// resolveDomains replaces each domain in t with its base, keeping the
// nullability t says each position has.
func resolveDomains(req *plugin.GenerateRequest, t *plugin.TypeExpr) (*plugin.TypeExpr, bool) {
	if d := findDomain(req, "", t.GetName()); !baseless(req, d) {
		base := proto.Clone(d.Base).(*plugin.TypeExpr)
		base.Nullable = t.Nullable
		if resolved, ok := resolveDomains(req, base); ok {
			return resolved, true
		}
		return base, true
	}
	var out *plugin.TypeExpr
	for i, arg := range t.GetArgs() {
		inner, ok := resolveDomains(req, arg.GetType())
		if arg.GetType() == nil || !ok {
			continue
		}
		if out == nil {
			out = proto.Clone(t).(*plugin.TypeExpr)
		}
		out.Args[i].Value = &plugin.TypeArg_Type{Type: inner}
	}
	return out, out != nil
}

// baseless reports whether d is missing or has no other type to resolve to:
// public.text over text reports its base as text, which a bare lookup
// finds as the domain again rather than the built-in it is.
func baseless(req *plugin.GenerateRequest, d *plugin.Domain) bool {
	return d == nil || d.Base == nil || findDomain(req, "", d.Base.GetName()) == d
}

// findDomain looks a domain up by the name a column's type reports: a
// qualified name in its schema, and a bare one only in the namespaces a
// type is reported from unqualified, since one elsewhere is reported
// qualified and so is a different type.
func findDomain(req *plugin.GenerateRequest, schema, name string) *plugin.Domain {
	if name == "" {
		return nil
	}
	if schema == "" {
		if i := strings.LastIndex(name, "."); i >= 0 {
			schema, name = name[:i], name[i+1:]
		}
	}
	schemas := []string{schema}
	if schema == "" {
		schemas = []string{req.GetCatalog().GetDefaultSchema(), "pg_catalog", engineDefaultSchema[req.GetSettings().GetEngine()]}
	}
	for _, s := range req.GetCatalog().GetSchemas() {
		if !slices.Contains(schemas, s.Name) {
			continue
		}
		for _, d := range s.GetDomains() {
			if d.Name == name {
				return d
			}
		}
	}
	return nil
}

// engineDefaultSchema is the schema an engine puts an unqualified object
// in when it is not the catalog's default, as its dialect seed names it.
var engineDefaultSchema = map[string]string{
	engineClickHouse: "default",
	engineDuckDB:     "main",
	engineMSSQL:      "dbo",
}
