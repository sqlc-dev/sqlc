package strongdb

import (
	q "github.com/kyleconroy/strongdb/tables/endpoint"
	t "github.com/kyleconroy/strongdb/tables/endpoint"
)

var scopedByAccount = q.Select{
	Where: q.And{
		q.Eq{t.AccountID, q.Arg},
		q.Eq{t.ID, q.Arg},
	},
}

var create = q.Insert{
	Returning: q.Star,
}

var update = q.Update{
	Set:       q.Columns{t.Settings},
	Where:     q.Eq{t.ID, q.Arg},
	Returning: q.Star,
}
