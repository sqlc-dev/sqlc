// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package datatype

import (
	"time"

	"github.com/jackc/pgtype"
)

type DtCharacter struct {
	A *string
	B *string
	C *string
	D *string
	E *string
}

type DtCharacterNotNull struct {
	A string
	B string
	C string
	D string
	E string
}

type DtDatetime struct {
	A *time.Time
	B *time.Time
	C *time.Time
	D *time.Time
	E *time.Time
	F *time.Time
	G *time.Time
	H *time.Time
}

type DtDatetimeNotNull struct {
	A time.Time
	B time.Time
	C time.Time
	D time.Time
	E time.Time
	F time.Time
	G time.Time
	H time.Time
}

type DtNetType struct {
	A pgtype.Inet
	B pgtype.CIDR
	C pgtype.Macaddr
}

type DtNetTypesNotNull struct {
	A pgtype.Inet
	B pgtype.CIDR
	C pgtype.Macaddr
}

type DtNumeric struct {
	A *int16
	B *int32
	C *int64
	D pgtype.Numeric
	E pgtype.Numeric
	F *float32
	G *float64
	H *int16
	I *int32
	J *int64
	K *int16
	L *int32
	M *int64
}

type DtNumericNotNull struct {
	A int16
	B int32
	C int64
	D pgtype.Numeric
	E pgtype.Numeric
	F float32
	G float64
	H int16
	I int32
	J int64
	K int16
	L int32
	M int64
}

type DtRange struct {
	A pgtype.Int4range
	B pgtype.Int8range
	C pgtype.Numrange
	D pgtype.Tsrange
	E pgtype.Tstzrange
	F pgtype.Daterange
}

type DtRangeNotNull struct {
	A pgtype.Int4range
	B pgtype.Int8range
	C pgtype.Numrange
	D pgtype.Tsrange
	E pgtype.Tstzrange
	F pgtype.Daterange
}
