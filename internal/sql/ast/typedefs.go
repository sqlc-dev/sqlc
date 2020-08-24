package ast

type AclMode uint32

func (n *AclMode) Pos() int {
	return 0
}

type DistinctExpr OpExpr

func (n *DistinctExpr) Pos() int {
	return 0
}

type NullIfExpr OpExpr

func (n *NullIfExpr) Pos() int {
	return 0
}

type Selectivity float64

func (n *Selectivity) Pos() int {
	return 0
}

type Cost float64

func (n *Cost) Pos() int {
	return 0
}

type ParamListInfo ParamListInfoData

func (n *ParamListInfo) Pos() int {
	return 0
}

type AttrNumber int16

func (n *AttrNumber) Pos() int {
	return 0
}

type Pointer byte

func (n *Pointer) Pos() int {
	return 0
}

type Index uint64

func (n *Index) Pos() int {
	return 0
}

type Offset int64

func (n *Offset) Pos() int {
	return 0
}

type regproc Oid

func (n *regproc) Pos() int {
	return 0
}

type RegProcedure regproc

func (n *RegProcedure) Pos() int {
	return 0
}

type TransactionId uint32

func (n *TransactionId) Pos() int {
	return 0
}

type LocalTransactionId uint32

func (n *LocalTransactionId) Pos() int {
	return 0
}

type SubTransactionId uint32

func (n *SubTransactionId) Pos() int {
	return 0
}

type MultiXactId TransactionId

func (n *MultiXactId) Pos() int {
	return 0
}

type MultiXactOffset uint32

func (n *MultiXactOffset) Pos() int {
	return 0
}

type CommandId uint32

func (n *CommandId) Pos() int {
	return 0
}

type Datum uintptr

func (n *Datum) Pos() int {
	return 0
}

type DatumPtr Datum

func (n *DatumPtr) Pos() int {
	return 0
}

type Oid uint64

func (n *Oid) Pos() int {
	return 0
}

type BlockNumber uint32

func (n *BlockNumber) Pos() int {
	return 0
}

type BlockId BlockIdData

func (n *BlockId) Pos() int {
	return 0
}
