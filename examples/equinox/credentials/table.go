package credentials

// CREATE TABLE credentials (
//         id         SERIAL       UNIQUE NOT NULL,
//         sid        varchar(64)  UNIQUE NOT NULL,
//         created    timestamp    DEFAULT NOW(),
//         accountid  bigint       NOT NULL,
//         tokenhash  varchar(255) NOT NULL
// )

type Table struct {
}

type Comparable interface {
}

type Column int

const (
	ID Column = iota
	SID
	Created
	AccountID
	TokenHash

	Arg
)

var Star = []Column{
	ID,
	SID,
	Created,
	AccountID,
	TokenHash,
}

type Select struct {
	Columns []Column
	Where   Comparable
}

type Eq struct {
	First  Column
	Second Column
}
