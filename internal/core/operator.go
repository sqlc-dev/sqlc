package core

import "fmt"

// OperatorSpec describes an operator overload for insertion into the catalog.
// A NULL left_type_oid encodes a prefix unary operator; NULL right_type_oid
// encodes a postfix unary.
type OperatorSpec struct {
	Name          string
	NamespaceOID  int64 // 0 = NULL
	DialectOID    int64 // 0 = NULL (shared)
	LeftTypeOID   int64 // 0 = NULL
	RightTypeOID  int64 // 0 = NULL
	ResultTypeOID int64
	ProcOID       int64 // 0 = NULL
}

// CreateOperator inserts an operator overload and returns its OID.
func (c *Catalog) CreateOperator(o OperatorSpec) (int64, error) {
	res, err := c.db.Exec(
		`INSERT INTO sql_operator
		   (namespace_oid, dialect_oid, name,
		    left_type_oid, right_type_oid, result_type_oid, proc_oid)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		nullableOID(o.NamespaceOID), nullableOID(o.DialectOID), o.Name,
		nullableOID(o.LeftTypeOID), nullableOID(o.RightTypeOID),
		o.ResultTypeOID, nullableOID(o.ProcOID),
	)
	if err != nil {
		return 0, fmt.Errorf("create operator %q: %w", o.Name, err)
	}
	return res.LastInsertId()
}

// OperatorOverload is a resolved candidate from FindOperators.
type OperatorOverload struct {
	OID           int64
	Name          string
	LeftTypeOID   int64 // 0 = NULL (prefix unary)
	RightTypeOID  int64 // 0 = NULL (postfix unary)
	ResultTypeOID int64
	ProcOID       int64 // 0 = NULL (binary-compatible)
}

// FindOperators returns all operator overloads matching the given name and
// (left,right) operand types. A 0 type means "any" and skips the filter on
// that side; useful for listing all overloads of a name.
func (c *Catalog) FindOperators(name string, leftTypeOID, rightTypeOID int64) ([]OperatorOverload, error) {
	q := `SELECT oid, name,
	             COALESCE(left_type_oid, 0), COALESCE(right_type_oid, 0),
	             result_type_oid, COALESCE(proc_oid, 0)
	      FROM sql_operator WHERE name = ?`
	args := []any{name}
	if leftTypeOID != 0 {
		q += ` AND left_type_oid = ?`
		args = append(args, leftTypeOID)
	}
	if rightTypeOID != 0 {
		q += ` AND right_type_oid = ?`
		args = append(args, rightTypeOID)
	}
	rows, err := c.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("find operators %q: %w", name, err)
	}
	defer rows.Close()
	var out []OperatorOverload
	for rows.Next() {
		var o OperatorOverload
		if err := rows.Scan(&o.OID, &o.Name, &o.LeftTypeOID, &o.RightTypeOID, &o.ResultTypeOID, &o.ProcOID); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
