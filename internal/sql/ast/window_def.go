package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type WindowDef struct {
	Name            *string
	Refname         *string
	PartitionClause *List
	OrderClause     *List
	FrameOptions    int
	StartOffset     Node
	EndOffset       Node
	Location        int
}

func (n *WindowDef) Pos() int {
	return n.Location
}

// Frame option constants (from PostgreSQL's parsenodes.h)
const (
	FrameOptionNonDefault              = 0x00001
	FrameOptionRange                   = 0x00002
	FrameOptionRows                    = 0x00004
	FrameOptionGroups                  = 0x00008
	FrameOptionBetween                 = 0x00010
	FrameOptionStartUnboundedPreceding = 0x00020
	FrameOptionEndUnboundedPreceding   = 0x00040
	FrameOptionStartUnboundedFollowing = 0x00080
	FrameOptionEndUnboundedFollowing   = 0x00100
	FrameOptionStartCurrentRow         = 0x00200
	FrameOptionEndCurrentRow           = 0x00400
	FrameOptionStartOffset             = 0x00800
	FrameOptionEndOffset               = 0x01000
	FrameOptionExcludeCurrentRow       = 0x02000
	FrameOptionExcludeGroup            = 0x04000
	FrameOptionExcludeTies             = 0x08000
)

func (n *WindowDef) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	// Named window reference
	if n.Refname != nil && *n.Refname != "" {
		buf.WriteString(*n.Refname)
		return
	}

	buf.WriteString("(")
	needSpace := false

	if items(n.PartitionClause) {
		buf.WriteString("PARTITION BY ")
		buf.join(n.PartitionClause, d, ", ")
		needSpace = true
	}

	if items(n.OrderClause) {
		if needSpace {
			buf.WriteString(" ")
		}
		buf.WriteString("ORDER BY ")
		buf.join(n.OrderClause, d, ", ")
		needSpace = true
	}

	// Frame clause
	if n.FrameOptions&FrameOptionNonDefault != 0 {
		if needSpace {
			buf.WriteString(" ")
		}

		// Frame type
		if n.FrameOptions&FrameOptionRows != 0 {
			buf.WriteString("ROWS ")
		} else if n.FrameOptions&FrameOptionRange != 0 {
			buf.WriteString("RANGE ")
		} else if n.FrameOptions&FrameOptionGroups != 0 {
			buf.WriteString("GROUPS ")
		}

		if n.FrameOptions&FrameOptionBetween != 0 {
			buf.WriteString("BETWEEN ")
		}

		// Start bound
		if n.FrameOptions&FrameOptionStartUnboundedPreceding != 0 {
			buf.WriteString("UNBOUNDED PRECEDING")
		} else if n.FrameOptions&FrameOptionStartCurrentRow != 0 {
			buf.WriteString("CURRENT ROW")
		} else if n.FrameOptions&FrameOptionStartOffset != 0 {
			buf.astFormat(n.StartOffset, d)
			buf.WriteString(" PRECEDING")
		}

		if n.FrameOptions&FrameOptionBetween != 0 {
			buf.WriteString(" AND ")

			// End bound
			if n.FrameOptions&FrameOptionEndUnboundedFollowing != 0 {
				buf.WriteString("UNBOUNDED FOLLOWING")
			} else if n.FrameOptions&FrameOptionEndCurrentRow != 0 {
				buf.WriteString("CURRENT ROW")
			} else if n.FrameOptions&FrameOptionEndOffset != 0 {
				buf.astFormat(n.EndOffset, d)
				buf.WriteString(" FOLLOWING")
			}
		}
	}

	buf.WriteString(")")
}
