// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

type Bar struct {
	ID    uint64
	Owner string
}

type Foo struct {
	Barid uint64
}
