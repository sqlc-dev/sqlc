// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package override

import (
	"github.com/Dionid/sqlc-testdata/pkg"
)

type Foo struct {
	Other   string
	Total   int64
	Retyped pkg.CustomType
}
