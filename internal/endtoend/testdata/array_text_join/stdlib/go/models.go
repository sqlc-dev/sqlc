// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.1

package querytest

import ()

type Bar struct {
	ID   string
	Info []string
}

type Foo struct {
	ID  string
	Bar string
}
