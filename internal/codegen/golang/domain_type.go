package golang

import (
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
		out.Type = &plugin.Identifier{Name: t.GetName()}
		return out
	}
	d := findDomain(req, col.Type.GetSchema(), col.Type.GetName())
	if d == nil || d.Base == nil {
		return nil
	}
	out := proto.Clone(col).(*plugin.Column)
	out.Type = &plugin.Identifier{Name: d.Base.Name}
	return out
}

// resolveDomains replaces each domain in t with its base, keeping the
// nullability t says each position has.
func resolveDomains(req *plugin.GenerateRequest, t *plugin.TypeExpr) (*plugin.TypeExpr, bool) {
	if d := findDomain(req, "", t.GetName()); d != nil && d.Base != nil {
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

// findDomain looks a domain up by the name a column's type reports: a
// qualified name in its schema, and a bare one in the default schema
// first, as the type would resolve.
func findDomain(req *plugin.GenerateRequest, schema, name string) *plugin.Domain {
	if name == "" {
		return nil
	}
	if schema == "" {
		if i := strings.LastIndex(name, "."); i >= 0 {
			schema, name = name[:i], name[i+1:]
		}
	}
	var fallback *plugin.Domain
	for _, s := range req.GetCatalog().GetSchemas() {
		for _, d := range s.GetDomains() {
			if d.Name != name {
				continue
			}
			switch {
			case schema != "" && s.Name == schema:
				return d
			case schema == "" && s.Name == req.Catalog.DefaultSchema:
				return d
			case schema == "" && fallback == nil:
				fallback = d
			}
		}
	}
	return fallback
}
