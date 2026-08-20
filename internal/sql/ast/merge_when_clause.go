package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

type MergeMatchKind uint

// Values match the PostgreSQL MergeMatchKind enum
const (
	MergeWhenMatched            MergeMatchKind = 1
	MergeWhenNotMatchedBySource MergeMatchKind = 2
	MergeWhenNotMatchedByTarget MergeMatchKind = 3
)

type MergeWhenClause struct {
	MatchKind   MergeMatchKind
	CommandType CmdType
	Override    OverridingKind
	Condition   Node
	TargetList  *List
	Values      *List
}

func (n *MergeWhenClause) Pos() int {
	return 0
}

func (n *MergeWhenClause) Format(buf *TrackedBuffer, d format.Dialect) {
	if n == nil {
		return
	}

	switch n.MatchKind {
	case MergeWhenNotMatchedBySource:
		buf.WriteString("WHEN NOT MATCHED BY SOURCE")
	case MergeWhenNotMatchedByTarget:
		buf.WriteString("WHEN NOT MATCHED")
	case MergeWhenMatched:
		buf.WriteString("WHEN MATCHED")
	}

	if set(n.Condition) {
		buf.WriteString(" AND ")
		buf.astFormat(n.Condition, d)
	}

	buf.WriteString(" THEN ")

	switch n.CommandType {
	case CmdTypeUpdate:
		buf.WriteString("UPDATE SET ")
		formatSetClause(buf, d, n.TargetList)
	case CmdTypeInsert:
		buf.WriteString("INSERT")
		if items(n.TargetList) {
			buf.WriteString(" (")
			buf.astFormat(n.TargetList, d)
			buf.WriteString(")")
		}
		switch n.Override {
		case OverridingSystemValue:
			buf.WriteString(" OVERRIDING SYSTEM VALUE")
		case OverridingUserValue:
			buf.WriteString(" OVERRIDING USER VALUE")
		}
		if items(n.Values) {
			buf.WriteString(" VALUES (")
			buf.join(n.Values, d, ", ")
			buf.WriteString(")")
		} else {
			buf.WriteString(" DEFAULT VALUES")
		}
	case CmdTypeDelete:
		buf.WriteString("DELETE")
	case CmdTypeNothing:
		buf.WriteString("DO NOTHING")
	}
}
