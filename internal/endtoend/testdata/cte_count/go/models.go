// Code generated by sqlc. DO NOT EDIT.

package querytest

import ()

type Bar struct {
	Ready bool
}

func (t *Bar) GetReady() bool {
	return t.Ready
}
