// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package override

import (
	orm "database/sql"
	"github.com/gofrs/uuid"
	fuid "github.com/gofrs/uuid"
	null "github.com/volatiletech/null/v8"
	null_v4 "gopkg.in/guregu/null.v4"
)

type Foo struct {
	ID      uuid.UUID
	OtherID fuid.UUID
	Age     orm.NullInt32
	Balance null.Float32
	Bio     null_v4.String
	About   *string
}
