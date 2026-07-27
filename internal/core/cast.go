package core

import (
	"context"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core/catalogdb"
)

// CastSpec describes a type-coercion rule.
//
// Context:
//   - 'i' implicit  — applied automatically by the analyzer
//   - 'a' assignment — applied for INSERT/UPDATE assignments
//   - 'e' explicit  — only when the user writes CAST or ::
//
// ProcOID == 0 means binary-coercible (no function call required).
type CastSpec struct {
	SourceTypeOID int64
	TargetTypeOID int64
	ProcOID       int64
	Context       string // 'i' | 'a' | 'e'; default 'e'
	DialectOID    int64
}

// CreateCast inserts a cast rule. (source, target) is the primary key, so
// re-inserting the same pair will fail; use ReplaceCast to overwrite.
func (c *Catalog) CreateCast(cs CastSpec) error {
	if cs.Context == "" {
		cs.Context = "e"
	}
	err := c.q.CreateCast(context.Background(), catalogdb.CreateCastParams{
		SourceTypeOid: cs.SourceTypeOID,
		TargetTypeOid: cs.TargetTypeOID,
		ProcOid:       nullInt64(cs.ProcOID),
		Context:       cs.Context,
		DialectOid:    nullInt64(cs.DialectOID),
	})
	if err != nil {
		return fmt.Errorf("create cast %d->%d: %w", cs.SourceTypeOID, cs.TargetTypeOID, err)
	}
	return nil
}

// FindCast returns the cast rule from src to tgt, if one exists, plus
// whether it was found.
func (c *Catalog) FindCast(src, tgt int64) (CastSpec, bool, error) {
	row, err := c.q.FindCast(context.Background(), catalogdb.FindCastParams{
		SourceTypeOid: src,
		TargetTypeOid: tgt,
	})
	if err != nil {
		return CastSpec{}, false, nil
	}
	return CastSpec{
		SourceTypeOID: row.SourceTypeOid,
		TargetTypeOID: row.TargetTypeOid,
		ProcOID:       orZero(row.ProcOid),
		Context:       row.Context,
		DialectOID:    orZero(row.DialectOid),
	}, true, nil
}

// CastAllowed reports whether src can be coerced to tgt in the given context.
// Same-type pairs always succeed.
func (c *Catalog) CastAllowed(src, tgt int64, ctx string) (bool, error) {
	if src == tgt {
		return true, nil
	}
	cs, ok, err := c.FindCast(src, tgt)
	if err != nil || !ok {
		return false, err
	}
	switch ctx {
	case "i":
		return cs.Context == "i", nil
	case "a":
		return cs.Context == "i" || cs.Context == "a", nil
	case "e":
		return true, nil
	}
	return false, fmt.Errorf("unknown cast context %q", ctx)
}
