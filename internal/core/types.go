package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// TypeSpec describes a type row. Name is the family's name; Expr is the
// canonical spelling of the whole expression and defaults to the name, which
// is what a family's is.
type TypeSpec struct {
	Name         string
	Expr         string
	Typtype      string
	Category     string
	Preferred    bool
	NamespaceOID int64
	DialectOID   int64
	FamilyOID    int64
	ElementOID   int64
	BaseOID      int64
	CanonicalOID int64
	NotNull      bool
}

// TypeInfo is a type row as the catalog holds it.
type TypeInfo struct {
	OID          int64
	NamespaceOID int64
	Name         string
	Expr         string
	Category     string
	Typtype      string
	Preferred    bool
	FamilyOID    int64
	ElementOID   int64
	BaseOID      int64
	CanonicalOID int64
	NotNull      bool
}

// IsFamily reports whether the row is a family rather than an instance.
func (t TypeInfo) IsFamily() bool { return t.FamilyOID == 0 }

// typeCache remembers what the catalog holds about a type. A row never
// changes once written, and a restored catalog is read-only, so a cached
// answer is good for the life of the catalog. Analysis runs concurrently on
// a restored catalog, so the cache is locked.
type typeCache struct {
	mu         sync.RWMutex
	infos      map[int64]TypeInfo
	exprs      map[int64]*TypeExpr
	namespaces map[int64]string
}

func (c *typeCache) namespace(oid int64) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	name, ok := c.namespaces[oid]
	return name, ok
}

func (c *typeCache) putNamespace(oid int64, name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.namespaces == nil {
		c.namespaces = map[int64]string{}
	}
	c.namespaces[oid] = name
}

func (c *typeCache) info(oid int64) (TypeInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	info, ok := c.infos[oid]
	return info, ok
}

func (c *typeCache) expr(oid int64) (*TypeExpr, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.exprs[oid]
	return e, ok
}

func (c *typeCache) put(info TypeInfo, expr *TypeExpr) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.infos == nil {
		c.infos = map[int64]TypeInfo{}
		c.exprs = map[int64]*TypeExpr{}
	}
	c.infos[info.OID] = info
	if expr != nil {
		c.exprs[info.OID] = expr
	}
}

func (c *Catalog) CreateType(name string) (int64, error) {
	return c.CreateTypeSpec(TypeSpec{Name: name, Typtype: "b"})
}

// CreateTypeSpec inserts a type row. It is the raw insert: nothing is
// canonicalized and no arguments are written, so it is what the seed and
// ResolveTypeExpr build on rather than what a caller with an expression
// wants.
func (c *Catalog) CreateTypeSpec(t TypeSpec) (int64, error) {
	if t.Typtype == "" {
		t.Typtype = "b"
	}
	if t.NamespaceOID == 0 {
		oid, err := c.NamespaceOID("public")
		if err != nil {
			return 0, fmt.Errorf("create type %q: default namespace: %w", t.Name, err)
		}
		t.NamespaceOID = oid
	}
	name := strings.ToLower(t.Name)
	expr := t.Expr
	if expr == "" {
		expr = name
	}
	oid, err := c.q.CreateType(context.Background(), catalogdb.CreateTypeParams{
		Name:         name,
		Expr:         expr,
		Typtype:      t.Typtype,
		Category:     nullString(t.Category),
		Preferred:    boolToInt64(t.Preferred),
		NamespaceOid: t.NamespaceOID,
		DialectOid:   nullInt64(t.DialectOID),
		FamilyOid:    nullInt64(t.FamilyOID),
		ElementOid:   nullInt64(t.ElementOID),
		BaseOid:      nullInt64(t.BaseOID),
		CanonicalOid: nullInt64(t.CanonicalOID),
		NotNull:      boolToInt64(t.NotNull),
	})
	if err != nil {
		return 0, fmt.Errorf("create type %q: %w", expr, err)
	}
	return oid, nil
}

// ArraySuffix is the suffix a schema appends to an element type's spelling
// to name an array of it, which ParseTypeExpr reads as one array dimension.
const ArraySuffix = "[]"

// CreateUserType registers a type family a schema declared rather than the
// dialect, such as an enum or a name the dialect's seed does not list. The
// type gains the dialect's comparison operators, so a column of it can be
// compared.
func (c *Catalog) CreateUserType(name, category string) (int64, error) {
	typtype := "b"
	if category == "E" {
		typtype = "e"
	}
	nsOID, bare, err := c.declaredTypeNamespace(strings.ToLower(name))
	if err != nil {
		return 0, fmt.Errorf("create type %q: %w", name, err)
	}
	spec := TypeSpec{
		Name:         bare,
		NamespaceOID: nsOID,
		Typtype:      typtype,
		Category:     category,
		DialectOID:   c.dialectOID,
	}
	// A dialect may say what an unknown spelling stands on, as SQLite's
	// affinity rule does; the type then resolves through that base and
	// needs no operators of its own.
	if category == "U" {
		baseOID, err := c.userTypeBase(bare)
		if err != nil {
			return 0, err
		}
		if baseOID != 0 {
			base, err := c.LookupType(baseOID)
			if err != nil {
				return 0, err
			}
			spec.BaseOID = baseOID
			spec.Category = base.Category
			return c.CreateTypeSpec(spec)
		}
	}
	oid, err := c.CreateTypeSpec(spec)
	if err != nil {
		return 0, err
	}
	if err := c.createComparisons(oid); err != nil {
		return 0, err
	}
	return oid, nil
}

// createComparisons gives a type the dialect's comparison operators, which the
// seed registered for the types it knew about up front.
func (c *Catalog) createComparisons(typeOID int64) error {
	if c.dialectOID == 0 {
		return nil
	}
	ops, _ := c.DialectFlag(c.dialectOID, FlagComparisonOperators)
	boolName, _ := c.DialectFlag(c.dialectOID, "const."+ConstBool)
	if ops == "" || boolName == "" {
		return nil
	}
	boolOID, err := c.TypeOID(boolName)
	if err != nil {
		return nil
	}
	for _, op := range strings.Split(ops, ",") {
		if op == "" {
			continue
		}
		if _, err := c.CreateOperator(OperatorSpec{
			Name:          op,
			DialectOID:    c.dialectOID,
			LeftTypeOID:   typeOID,
			RightTypeOID:  typeOID,
			ResultTypeOID: boolOID,
		}); err != nil {
			return err
		}
	}
	return nil
}

// TypeOIDsInCategory returns the type families the catalog's dialect has in
// the named category, in the order they were created.
func (c *Catalog) TypeOIDsInCategory(category string) ([]int64, error) {
	oids, err := c.q.TypeOIDsInCategory(context.Background(), catalogdb.TypeOIDsInCategoryParams{
		DialectOid: nullInt64(c.dialectOID),
		Category:   nullString(category),
	})
	if err != nil {
		return nil, fmt.Errorf("types in category %q: %w", category, err)
	}
	return oids, nil
}

// TypeOID returns the family a name refers to: an alias spelling resolves to
// the type it names.
func (c *Catalog) TypeOID(name string) (int64, error) {
	oid, err := c.familyOIDByName(strings.ToLower(name))
	if err != nil {
		return 0, fmt.Errorf("type %q: %w", name, err)
	}
	return c.canonicalOID(oid)
}

// familyOIDByName finds the family row spelled name, alias rows included.
func (c *Catalog) familyOIDByName(name string) (int64, error) {
	return c.q.TypeOIDByName(context.Background(), name)
}

// familyOIDByQualifiedName is familyOIDByName for a name that may carry its
// namespace, as myschema.mood does: a qualified name is looked up in that
// namespace alone, a bare one in every namespace.
func (c *Catalog) familyOIDByQualifiedName(name string) (int64, error) {
	ns, bare := splitQualifiedName(name)
	if ns == "" {
		return c.familyOIDByName(bare)
	}
	nsOID, err := c.NamespaceOID(ns)
	if err != nil {
		return 0, err
	}
	return c.q.TypeOIDByNameInNamespace(context.Background(), catalogdb.TypeOIDByNameInNamespaceParams{
		NamespaceOid: nsOID,
		Name:         bare,
	})
}

// splitQualifiedName splits "myschema.mood" into its namespace and name. A
// name with no dot has no namespace.
func splitQualifiedName(name string) (ns, bare string) {
	if i := strings.LastIndexByte(name, '.'); i > 0 {
		return name[:i], name[i+1:]
	}
	return "", name
}

// declaredTypeNamespace is the namespace a declared type's row goes in: the
// one its name qualifies, created if the schema has not, or the default.
func (c *Catalog) declaredTypeNamespace(name string) (int64, string, error) {
	ns, bare := splitQualifiedName(name)
	if ns == "" {
		return 0, bare, nil
	}
	oid, err := c.NamespaceOID(ns)
	if err != nil {
		if oid, err = c.CreateNamespace(ns); err != nil {
			return 0, "", err
		}
	}
	return oid, bare, nil
}

// CreateTypeWithArgs registers a declared type that has arguments of its own
// — a composite's fields, an enum's labels — as a family row carrying them.
// The arguments' types are interned first.
func (c *Catalog) CreateTypeWithArgs(spec TypeSpec, args []TypeArg) (int64, error) {
	nsOID, bare, err := c.declaredTypeNamespace(strings.ToLower(spec.Name))
	if err != nil {
		return 0, fmt.Errorf("create type %q: %w", spec.Name, err)
	}
	spec.Name = bare
	if nsOID != 0 {
		spec.NamespaceOID = nsOID
	}
	if spec.DialectOID == 0 {
		spec.DialectOID = c.dialectOID
	}
	argOIDs := make([]int64, len(args))
	for i, a := range args {
		if a.Type == nil {
			continue
		}
		oid, _, err := c.internType(a.Type, func(name string) (int64, error) {
			return c.CreateUserType(name, "U")
		})
		if err != nil {
			return 0, fmt.Errorf("create type %q: %w", spec.Name, err)
		}
		argOIDs[i] = oid
	}
	oid, err := c.CreateTypeSpec(spec)
	if err != nil {
		return 0, err
	}
	if err := c.createComparisons(oid); err != nil {
		return 0, err
	}
	return oid, c.insertTypeArgs(oid, spec.Expr, args, argOIDs)
}

// insertTypeArgs writes a type's argument rows.
func (c *Catalog) insertTypeArgs(oid int64, key string, args []TypeArg, argOIDs []int64) error {
	ctx := context.Background()
	for i, a := range args {
		p := catalogdb.CreateTypeArgParams{TypeOid: oid, Ord: int64(i + 1), Label: a.Label}
		switch {
		case a.Type != nil:
			p.ArgTypeOid = nullInt64(argOIDs[i])
			p.Nullable = boolToInt64(a.Type.Nullable)
		case a.Int != nil:
			p.IntValue = sql.NullInt64{Int64: *a.Int, Valid: true}
		case a.Bool != nil:
			p.BoolValue = sql.NullInt64{Int64: boolToInt64(*a.Bool), Valid: true}
		case a.String != nil:
			p.StringValue = sql.NullString{String: *a.String, Valid: true}
		case a.Ident != nil:
			p.Ident = sql.NullString{String: *a.Ident, Valid: true}
		}
		if err := c.q.CreateTypeArg(ctx, p); err != nil {
			return fmt.Errorf("type %q: argument %d: %w", key, i+1, err)
		}
	}
	return nil
}

// canonicalOID follows an alias row to the row it stands for.
func (c *Catalog) canonicalOID(oid int64) (int64, error) {
	for i := 0; i < 16; i++ {
		info, err := c.LookupType(oid)
		if err != nil {
			return 0, err
		}
		if info.CanonicalOID == 0 {
			return oid, nil
		}
		oid = info.CanonicalOID
	}
	return 0, fmt.Errorf("type oid %d: alias chain does not end", oid)
}

// ResolutionOID is the row an operator, function or cast over the type is
// looked up on when none is registered on the type itself: an instance's
// family, an alias's canonical type, a domain's or wrapper's base. It
// returns 0 when there is nothing further to fall back to.
func (c *Catalog) ResolutionOID(oid int64) int64 {
	info, err := c.LookupType(oid)
	if err != nil {
		return 0
	}
	switch {
	case info.CanonicalOID != 0:
		return info.CanonicalOID
	case info.FamilyOID != 0:
		return info.FamilyOID
	case info.BaseOID != 0:
		return info.BaseOID
	}
	return 0
}

// ResolutionChain lists the type and every row resolution falls back to, in
// order, ending with the family everything about the type is registered on.
func (c *Catalog) ResolutionChain(oid int64) []int64 {
	chain := []int64{oid}
	for i := 0; i < 16 && oid != 0; i++ {
		oid = c.ResolutionOID(oid)
		if oid != 0 {
			chain = append(chain, oid)
		}
	}
	return chain
}

// TypeName returns the family name of a type row.
func (c *Catalog) TypeName(oid int64) (string, error) {
	info, err := c.LookupType(oid)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}

// LookupType returns what the catalog holds about a type row.
func (c *Catalog) LookupType(oid int64) (TypeInfo, error) {
	if info, ok := c.types.info(oid); ok {
		return info, nil
	}
	row, err := c.q.LookupType(context.Background(), oid)
	if err != nil {
		return TypeInfo{}, fmt.Errorf("lookup type oid %d: %w", oid, err)
	}
	info := TypeInfo{
		OID:          row.Oid,
		NamespaceOID: row.NamespaceOid,
		Name:         row.Name,
		Expr:         row.Expr,
		Category:     row.Category.String,
		Typtype:      row.Typtype,
		Preferred:    row.Preferred != 0,
		FamilyOID:    orZero(row.FamilyOid),
		ElementOID:   orZero(row.ElementOid),
		BaseOID:      orZero(row.BaseOid),
		CanonicalOID: orZero(row.CanonicalOid),
		NotNull:      row.NotNull != 0,
	}
	c.types.put(info, nil)
	return info, nil
}

// namespaceName is the name of a namespace row, remembered once read.
func (c *Catalog) namespaceName(oid int64) (string, error) {
	if name, ok := c.types.namespace(oid); ok {
		return name, nil
	}
	namespaces, err := c.Namespaces()
	if err != nil {
		return "", err
	}
	for _, ns := range namespaces {
		c.types.putNamespace(ns.OID, ns.Name)
	}
	name, _ := c.types.namespace(oid)
	return name, nil
}

// TypeExprOf is the expression a type row stands for, read back from its
// arguments: the family's name for a family, the family applied to its
// arguments for an instance. The result is the caller's to change.
func (c *Catalog) TypeExprOf(oid int64) (*TypeExpr, error) {
	if e, ok := c.types.expr(oid); ok {
		return e.Clone(), nil
	}
	info, err := c.LookupType(oid)
	if err != nil {
		return nil, err
	}
	expr := &TypeExpr{Name: info.Name}
	// A type outside the default namespaces is named with its namespace,
	// as format_type prints a type off the search path.
	if ns, err := c.namespaceName(info.NamespaceOID); err == nil && ns != "" && !slices.Contains(c.DefaultNamespaces(), ns) {
		expr.Name = ns + "." + info.Name
	}
	if !info.IsFamily() {
		rows, err := c.q.TypeArgs(context.Background(), oid)
		if err != nil {
			return nil, fmt.Errorf("type oid %d: arguments: %w", oid, err)
		}
		for _, r := range rows {
			arg := TypeArg{Label: r.Label}
			switch {
			case r.ArgTypeOid.Valid:
				t, err := c.TypeExprOf(r.ArgTypeOid.Int64)
				if err != nil {
					return nil, err
				}
				t.Nullable = r.Nullable != 0
				arg.Type = t
			case r.IntValue.Valid:
				v := r.IntValue.Int64
				arg.Int = &v
			case r.BoolValue.Valid:
				v := r.BoolValue.Int64 != 0
				arg.Bool = &v
			case r.StringValue.Valid:
				v := r.StringValue.String
				arg.String = &v
			case r.Ident.Valid:
				v := r.Ident.String
				arg.Ident = &v
			}
			expr.Args = append(expr.Args, arg)
		}
	}
	c.types.put(info, expr)
	return expr.Clone(), nil
}

// ResolveTypeExpr interns the type an expression names and returns its row:
// the family for a bare name, the instance for a family applied to
// arguments, each argument type interned first. Names are canonicalized, so
// integer and int4 intern to one row. A family the dialect did not seed is
// one the schema declared, and is registered as a user type.
func (c *Catalog) ResolveTypeExpr(t *TypeExpr) (int64, error) {
	oid, _, err := c.internType(t, func(name string) (int64, error) {
		return c.CreateUserType(name, "U")
	})
	return oid, err
}

// ResolvePseudoTypeExpr is ResolveTypeExpr for the type a function signature
// names. Signatures reference pseudo-types ("any", "record") and types no
// dialect bothers to list, so an unknown family is registered as an opaque
// type rather than rejected, and gets no operators of its own.
func (c *Catalog) ResolvePseudoTypeExpr(t *TypeExpr) (int64, error) {
	oid, _, err := c.internType(t, func(name string) (int64, error) {
		return c.CreateTypeSpec(TypeSpec{Name: name, Category: "U", DialectOID: c.dialectOID})
	})
	return oid, err
}

// ResolveTypeName is ResolveTypeExpr for a type spelled as a string.
func (c *Catalog) ResolveTypeName(name string) (int64, error) {
	return c.ResolveTypeExpr(ParseTypeExpr(name))
}

var errUnknownType = errors.New("unknown type")

// TypeLookup is what LookupTypeExpr found for an expression.
type TypeLookup struct {
	// OID is the instance row when the catalog holds one, otherwise the
	// family row.
	OID int64
	// FamilyOID is the family, which is OID for a family or an instance the
	// catalog does not hold.
	FamilyOID int64
	// Expr is the expression canonicalized: the family and every argument
	// type spelled as the catalog spells them, whether or not the instance
	// is a row.
	Expr *TypeExpr
}

// LookupTypeExpr finds the row an expression names without writing, and
// reports false when the family is not one the catalog holds.
func (c *Catalog) LookupTypeExpr(t *TypeExpr) (TypeLookup, bool) {
	refuse := func(string) (int64, error) { return 0, errUnknownType }
	oid, canonical, err := c.internType(t, refuse)
	if errors.Is(err, errUnknownType) && canonical != nil {
		// The family is known and the instance is not a row.
		familyOID, _, err := c.internType(&TypeExpr{Name: canonical.Name}, refuse)
		if err != nil {
			return TypeLookup{}, false
		}
		return TypeLookup{OID: familyOID, FamilyOID: familyOID, Expr: canonical}, true
	}
	if err != nil {
		return TypeLookup{}, false
	}
	info, err := c.LookupType(oid)
	if err != nil {
		return TypeLookup{}, false
	}
	familyOID := oid
	if !info.IsFamily() {
		familyOID = info.FamilyOID
	}
	return TypeLookup{OID: oid, FamilyOID: familyOID, Expr: canonical}, true
}

// internType resolves an expression to its row, creating the instance row
// when there is none and calling newFamily for a family name the catalog
// does not hold, which may refuse. Alongside the row it returns the
// expression canonicalized; when the instance is not a row and newFamily
// refuses, the canonical expression still comes back with the error.
func (c *Catalog) internType(t *TypeExpr, newFamily func(name string) (int64, error)) (int64, *TypeExpr, error) {
	if t == nil || strings.TrimSpace(t.Name) == "" {
		return 0, nil, fmt.Errorf("missing type name")
	}
	t, err := c.canonicalize(t)
	if err != nil {
		return 0, nil, err
	}
	name := strings.ToLower(strings.TrimSpace(t.Name))
	familyOID, err := c.familyOIDByQualifiedName(name)
	if err != nil {
		if name == ArrayTypeName {
			// Every dialect has arrays, whether or not its seed lists the
			// family; one that does not gets it as an array type rather
			// than a user type.
			familyOID, err = c.CreateTypeSpec(TypeSpec{Name: name, Category: "A", DialectOID: c.dialectOID})
		} else {
			familyOID, err = newFamily(name)
		}
		if err != nil {
			return 0, nil, fmt.Errorf("type %q: %w", name, err)
		}
	}
	if familyOID, err = c.canonicalOID(familyOID); err != nil {
		return 0, nil, err
	}
	family, err := c.LookupType(familyOID)
	if err != nil {
		return 0, nil, err
	}
	if len(t.Args) == 0 {
		return familyOID, &TypeExpr{Name: family.Name}, nil
	}

	// The instance's canonical spelling is the family applied to its
	// arguments as the catalog spells them, so each argument type is
	// resolved first and read back.
	canonical := &TypeExpr{Name: family.Name, Args: make([]TypeArg, len(t.Args))}
	argOIDs := make([]int64, len(t.Args))
	for i, a := range t.Args {
		canonical.Args[i] = a
		if a.Type == nil {
			continue
		}
		oid, argExpr, err := c.internType(a.Type, newFamily)
		if err != nil {
			return 0, nil, err
		}
		argOIDs[i] = oid
		argExpr.Nullable = a.Type.Nullable
		canonical.Args[i].Type = argExpr
	}
	key := canonical.Key()
	ctx := context.Background()
	if oid, err := c.q.TypeOIDByExprInNamespace(ctx, catalogdb.TypeOIDByExprInNamespaceParams{
		NamespaceOid: family.NamespaceOID,
		Expr:         key,
	}); err == nil {
		return oid, canonical, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, nil, fmt.Errorf("type %q: %w", key, err)
	}
	// A lookup that may not write stops here, canonical expression in hand.
	if _, err := newFamily(""); errors.Is(err, errUnknownType) {
		return 0, canonical, errUnknownType
	}

	spec := TypeSpec{
		Name:         family.Name,
		Expr:         key,
		Typtype:      family.Typtype,
		Category:     family.Category,
		NamespaceOID: family.NamespaceOID,
		DialectOID:   c.dialectOID,
		FamilyOID:    familyOID,
	}
	if family.Name == ArrayTypeName && argOIDs[0] != 0 {
		spec.ElementOID = argOIDs[0]
	}
	oid, err := c.CreateTypeSpec(spec)
	if err != nil {
		return 0, nil, err
	}
	if err := c.insertTypeArgs(oid, key, canonical.Args, argOIDs); err != nil {
		return 0, nil, err
	}
	return oid, canonical, nil
}

func nullableOID(oid int64) any {
	if oid == 0 {
		return nil
	}
	return oid
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullInt64(oid int64) sql.NullInt64 {
	return sql.NullInt64{Int64: oid, Valid: oid != 0}
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func orZero(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}
