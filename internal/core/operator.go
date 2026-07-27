package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

type OperatorSpec struct {
	Name          string
	NamespaceOID  int64
	DialectOID    int64
	LeftTypeOID   int64
	RightTypeOID  int64
	ResultTypeOID int64
	ProcOID       int64
}

func (c *Catalog) CreateOperator(o OperatorSpec) (int64, error) {
	oid, err := c.q.CreateOperator(context.Background(), catalogdb.CreateOperatorParams{
		NamespaceOid:  nullInt64(o.NamespaceOID),
		DialectOid:    nullInt64(o.DialectOID),
		Name:          o.Name,
		LeftTypeOid:   nullInt64(o.LeftTypeOID),
		RightTypeOid:  nullInt64(o.RightTypeOID),
		ResultTypeOid: o.ResultTypeOID,
		ProcOid:       nullInt64(o.ProcOID),
	})
	if err != nil {
		return 0, fmt.Errorf("create operator %q: %w", o.Name, err)
	}
	return oid, nil
}

type OperatorOverload struct {
	OID           int64
	Name          string
	LeftTypeOID   int64
	RightTypeOID  int64
	ResultTypeOID int64
	ProcOID       int64
}

func (c *Catalog) FindOperators(name string, leftTypeOID, rightTypeOID int64) ([]OperatorOverload, error) {
	rows, err := c.q.FindOperators(context.Background(), catalogdb.FindOperatorsParams{
		Name:         name,
		LeftTypeOid:  leftTypeOID,
		RightTypeOid: rightTypeOID,
	})
	if err != nil {
		return nil, fmt.Errorf("find operators %q: %w", name, err)
	}
	out := make([]OperatorOverload, 0, len(rows))
	for _, r := range rows {
		out = append(out, OperatorOverload{
			OID:           r.Oid,
			Name:          r.Name,
			LeftTypeOID:   orZero(r.LeftTypeOid),
			RightTypeOID:  orZero(r.RightTypeOid),
			ResultTypeOID: r.ResultTypeOid,
			ProcOID:       orZero(r.ProcOid),
		})
	}
	return out, nil
}
