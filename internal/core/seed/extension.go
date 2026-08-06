package seed

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
)

// applyExtension applies the named extension's directory to a catalog that
// has already been seeded. Unlike the dialect's own seed, an extension lands
// in a catalog full of types, so everything it names is resolved against what
// is there before being created.
func applyExtension(cat *core.Catalog, fsys fs.FS, name string) error {
	dir := path.Join(ExtensionsDir, name)
	if _, err := fs.Stat(fsys, dir); err != nil {
		// An extension sqlc has no data for adds nothing, the way the legacy
		// catalog has always treated one.
		return nil
	}
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return fmt.Errorf("seed: extension %s: %w", name, err)
	}
	e := &extension{cat: cat}
	if err := stream(sub, TypesFile, e.addType); err != nil {
		return fmt.Errorf("extension %s: %w", name, err)
	}
	if err := stream(sub, FunctionsFile, e.addFunction); err != nil {
		return fmt.Errorf("extension %s: %w", name, err)
	}
	return nil
}

type extension struct {
	cat *core.Catalog
}

// addType registers a type the extension defines. The type joins the
// dialect's rules the way a type declared by a schema does: it gains the
// comparison operators, and implicit casts to the types of its category when
// the dialect's seed declared that category mutually castable.
func (e *extension) addType(t Type) error {
	for _, name := range append([]string{t.Name}, t.Aliases...) {
		if _, err := e.cat.TypeOID(name); err == nil {
			// The dialect, or an earlier extension, already ships this type.
			continue
		}
		oid, err := e.cat.CreateUserType(name, t.Category)
		if err != nil {
			return fmt.Errorf("type %q: %w", name, err)
		}
		if err := e.categoryCasts(oid, t.Category); err != nil {
			return fmt.Errorf("type %q: %w", name, err)
		}
	}
	return nil
}

// categoryCasts makes a new type implicitly castable to and from the types
// already in its category, under the same rule the dialect's seed applied to
// the types it was born with.
func (e *extension) categoryCasts(oid int64, category string) error {
	categories, err := e.cat.DialectFlag(e.cat.SeededDialectOID(), core.FlagCastCategories)
	if err != nil || categories == "" {
		return err
	}
	if categories != "*" && !strings.Contains(categories, category) {
		return nil
	}
	peers, err := e.cat.TypeOIDsInCategory(category)
	if err != nil {
		return err
	}
	for _, peer := range peers {
		if peer == oid {
			continue
		}
		for _, cast := range [][2]int64{{oid, peer}, {peer, oid}} {
			if err := e.cat.CreateCast(core.CastSpec{
				SourceTypeOID: cast[0],
				TargetTypeOID: cast[1],
				Context:       "i",
				DialectOID:    e.cat.SeededDialectOID(),
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *extension) addFunction(fn Function) error {
	returnOID, err := e.funcType(fn.Returns)
	if err != nil {
		return fmt.Errorf("function %q: %w", fn.Name, err)
	}
	args := make([]core.ProcArg, 0, len(fn.Args))
	for _, a := range fn.Args {
		// An out or table parameter is part of the result, not the call.
		switch a.Mode {
		case "o", "t":
			continue
		}
		argOID, err := e.funcType(a.Type)
		if err != nil {
			return fmt.Errorf("function %q: %w", fn.Name, err)
		}
		args = append(args, core.ProcArg{
			Name:       a.Name,
			TypeOID:    argOID,
			Mode:       a.Mode,
			HasDefault: a.HasDefault,
		})
	}
	_, err = e.cat.CreateProc(core.ProcSpec{
		Name:           fn.Name,
		DialectOID:     e.cat.SeededDialectOID(),
		Kind:           fn.Kind,
		ReturnTypeOID:  returnOID,
		ReturnNullable: fn.Nullable,
		Args:           args,
	})
	if err != nil {
		return fmt.Errorf("function %q: %w", fn.Name, err)
	}
	return nil
}

// funcType resolves a type a function signature names, registering the ones
// no seed bothered to list the way the dialect's own function list does.
func (e *extension) funcType(name string) (int64, error) {
	if name == "" {
		return 0, nil
	}
	if oid, err := e.cat.TypeOID(name); err == nil {
		return oid, nil
	}
	return e.cat.CreateTypeSpec(core.TypeSpec{
		Name:       name,
		Typtype:    "b",
		Category:   "U",
		DialectOID: e.cat.SeededDialectOID(),
	})
}
