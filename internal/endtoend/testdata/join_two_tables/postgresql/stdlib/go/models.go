// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package querytest

type Bar struct {
	ID int32
}

type Baz struct {
	ID int32
}

type Foo struct {
	BarID int32
	BazID int32
}
