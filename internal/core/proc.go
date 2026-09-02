package core

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

type ProcSpec struct {
	Name           string
	NamespaceOID   int64
	DialectOID     int64
	Kind           string
	ReturnTypeOID  int64
	ReturnSet      bool
	ReturnNullable bool
	// NeverNull marks a function whose result is never NULL even when an
	// argument is, in a dialect that otherwise propagates nullability.
	NeverNull    bool
	Strict       bool
	VariadicKind string
	Args         []ProcArg
}

// The proc table stores nullability as one integer: 0 leaves it to the
// dialect's rule, 1 is always nullable and 2 is never nullable.
const (
	nullableDefault int64 = 0
	nullableAlways  int64 = 1
	nullableNever   int64 = 2
)

type ProcArg struct {
	Name       string
	TypeOID    int64
	Mode       string
	HasDefault bool
}

func (c *Catalog) CreateProc(p ProcSpec) (int64, error) {
	if p.Kind == "" {
		p.Kind = "f"
	}
	if p.VariadicKind == "" {
		p.VariadicKind = "n"
	}
	ctx := context.Background()
	procOID, err := c.q.CreateProc(ctx, catalogdb.CreateProcParams{
		NamespaceOid:   nullInt64(p.NamespaceOID),
		DialectOid:     nullInt64(p.DialectOID),
		Name:           strings.ToLower(p.Name),
		Kind:           p.Kind,
		ReturnTypeOid:  p.ReturnTypeOID,
		ReturnSet:      boolToInt64(p.ReturnSet),
		ReturnNullable: returnNullable(p),
		Strict:         boolToInt64(p.Strict),
		VariadicKind:   p.VariadicKind,
	})
	if err != nil {
		return 0, fmt.Errorf("create proc %q: %w", p.Name, err)
	}
	for i, a := range p.Args {
		mode := a.Mode
		if mode == "" {
			mode = "i"
		}
		err := c.q.CreateProcArg(ctx, catalogdb.CreateProcArgParams{
			ProcOid:    procOID,
			Ord:        int64(i + 1),
			Name:       a.Name,
			TypeOid:    a.TypeOID,
			Mode:       mode,
			HasDefault: boolToInt64(a.HasDefault),
		})
		if err != nil {
			return 0, fmt.Errorf("create proc %q arg %d: %w", p.Name, i+1, err)
		}
	}
	return procOID, nil
}

func returnNullable(p ProcSpec) int64 {
	switch {
	case p.NeverNull:
		return nullableNever
	case p.ReturnNullable:
		return nullableAlways
	}
	return nullableDefault
}

type ProcOverload struct {
	OID            int64
	Name           string
	Kind           string
	ReturnTypeOID  int64
	ReturnNullable bool
	NeverNull      bool
	ArgTypes       []int64
}

// FindProcs returns the overloads of name, optionally restricted to the given
// namespaces.
func (c *Catalog) FindProcs(name string, namespaceOIDs []int64) ([]ProcOverload, error) {
	ctx := context.Background()
	lname := strings.ToLower(name)

	var out []ProcOverload
	if len(namespaceOIDs) == 0 {
		rows, err := c.q.FindProcsAnyNamespace(ctx, lname)
		if err != nil {
			return nil, fmt.Errorf("find procs %q: %w", name, err)
		}
		out = make([]ProcOverload, 0, len(rows))
		for _, r := range rows {
			out = append(out, ProcOverload{
				OID:            r.Oid,
				Name:           r.Name,
				Kind:           r.Kind,
				ReturnTypeOID:  r.ReturnTypeOid,
				ReturnNullable: r.ReturnNullable == nullableAlways,
				NeverNull:      r.ReturnNullable == nullableNever,
			})
		}
	} else {
		nss := make([]sql.NullInt64, len(namespaceOIDs))
		for i, ns := range namespaceOIDs {
			nss[i] = sql.NullInt64{Int64: ns, Valid: true}
		}
		rows, err := c.q.FindProcsInNamespaces(ctx, catalogdb.FindProcsInNamespacesParams{
			Name:          lname,
			NamespaceOids: nss,
		})
		if err != nil {
			return nil, fmt.Errorf("find procs %q: %w", name, err)
		}
		out = make([]ProcOverload, 0, len(rows))
		for _, r := range rows {
			out = append(out, ProcOverload{
				OID:            r.Oid,
				Name:           r.Name,
				Kind:           r.Kind,
				ReturnTypeOID:  r.ReturnTypeOid,
				ReturnNullable: r.ReturnNullable == nullableAlways,
				NeverNull:      r.ReturnNullable == nullableNever,
			})
		}
	}

	for i := range out {
		argTypes, err := c.procArgTypes(out[i].OID)
		if err != nil {
			return nil, err
		}
		out[i].ArgTypes = argTypes
	}
	return out, nil
}

func (c *Catalog) procArgTypes(procOID int64) ([]int64, error) {
	return c.q.ProcArgTypes(context.Background(), procOID)
}
