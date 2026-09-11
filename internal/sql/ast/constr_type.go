package ast

type ConstrType uint

func (n *ConstrType) Pos() int {
	return 0
}

// The constraint kinds the analysis reads, numbered as PostgreSQL's parser
// numbers them, which is what an engine's converter records.
const (
	ConstrTypeNull    ConstrType = 1
	ConstrTypeNotNull ConstrType = 2
)
