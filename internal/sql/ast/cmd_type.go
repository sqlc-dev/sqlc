package ast

type CmdType uint

// Values match the PostgreSQL CmdType enum
const (
	CmdTypeUnknown CmdType = 1
	CmdTypeSelect  CmdType = 2
	CmdTypeUpdate  CmdType = 3
	CmdTypeInsert  CmdType = 4
	CmdTypeDelete  CmdType = 5
	CmdTypeMerge   CmdType = 6
	CmdTypeUtility CmdType = 7
	CmdTypeNothing CmdType = 8
)

func (n *CmdType) Pos() int {
	return 0
}
