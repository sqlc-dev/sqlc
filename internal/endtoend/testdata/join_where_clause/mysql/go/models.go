// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package querytest

import ()

type Bar struct {
	ID    uint64
	Owner string
}

type Foo struct {
	Barid uint64
}
