package core

import (
	"database/sql"
	"fmt"
	"strings"
)

// ProcSpec describes a function/aggregate/window/procedure for insertion
// into the catalog.
type ProcSpec struct {
	Name           string
	NamespaceOID   int64  // 0 = NULL
	DialectOID     int64  // 0 = NULL (shared)
	Kind           string // 'f'unc | 'a'gg | 'w'in | 'p'roc; default 'f'
	ReturnTypeOID  int64
	ReturnSet      bool
	ReturnNullable bool
	Strict         bool
	VariadicKind   string // 'n'one | 'a'rray | 'v'ariadic-any; default 'n'
	Args           []ProcArg
}

// ProcArg is a single argument in a proc signature.
type ProcArg struct {
	Name       string
	TypeOID    int64
	Mode       string // 'i'n | 'o'ut | 'b'oth | 't'able | 'v'ariadic; default 'i'
	HasDefault bool
}

// CreateProc inserts a proc and its arguments, returning the proc OID.
func (c *Catalog) CreateProc(p ProcSpec) (int64, error) {
	if p.Kind == "" {
		p.Kind = "f"
	}
	if p.VariadicKind == "" {
		p.VariadicKind = "n"
	}
	res, err := c.db.Exec(
		`INSERT INTO sql_proc
		   (namespace_oid, dialect_oid, name, kind,
		    return_type_oid, return_set, return_nullable, strict, variadic_kind)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		nullableOID(p.NamespaceOID), nullableOID(p.DialectOID),
		strings.ToLower(p.Name), p.Kind,
		p.ReturnTypeOID, boolToInt(p.ReturnSet), boolToInt(p.ReturnNullable),
		boolToInt(p.Strict), p.VariadicKind,
	)
	if err != nil {
		return 0, fmt.Errorf("create proc %q: %w", p.Name, err)
	}
	procOID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for i, a := range p.Args {
		mode := a.Mode
		if mode == "" {
			mode = "i"
		}
		_, err := c.db.Exec(
			`INSERT INTO sql_proc_arg (proc_oid, ord, name, type_oid, mode, has_default)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			procOID, i+1, a.Name, a.TypeOID, mode, boolToInt(a.HasDefault),
		)
		if err != nil {
			return 0, fmt.Errorf("create proc %q arg %d: %w", p.Name, i+1, err)
		}
	}
	return procOID, nil
}

// ProcOverload describes one resolved candidate from FindProcs.
type ProcOverload struct {
	OID            int64
	Name           string
	Kind           string
	ReturnTypeOID  int64
	ReturnNullable bool
	ArgTypes       []int64
}

// FindProcs returns all procs matching the given (case-insensitive) name in
// any of the supplied namespaces. Pass an empty namespace list to search
// every namespace. Argument lists are populated in declaration order.
func (c *Catalog) FindProcs(name string, namespaceOIDs []int64) ([]ProcOverload, error) {
	q := `SELECT oid, name, kind, return_type_oid, return_nullable
	      FROM sql_proc WHERE name = ?`
	args := []any{strings.ToLower(name)}
	if len(namespaceOIDs) > 0 {
		q += " AND namespace_oid IN (" + placeholders(len(namespaceOIDs)) + ")"
		for _, ns := range namespaceOIDs {
			args = append(args, ns)
		}
	}
	rows, err := c.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("find procs %q: %w", name, err)
	}
	defer rows.Close()
	var out []ProcOverload
	for rows.Next() {
		var o ProcOverload
		var rn int
		if err := rows.Scan(&o.OID, &o.Name, &o.Kind, &o.ReturnTypeOID, &rn); err != nil {
			return nil, err
		}
		o.ReturnNullable = rn != 0
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
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
	rows, err := c.db.Query(
		`SELECT type_oid FROM sql_proc_arg
		 WHERE proc_oid = ? AND mode IN ('i','b','v')
		 ORDER BY ord`, procOID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var t int64
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

// silence unused-import warning if we drop the sql.Null* usages
var _ sql.NullString
