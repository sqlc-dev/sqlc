package strongdb

import (
	"context"
	"database/sql"
	"encoding/json"
)

type EndpointConfig struct {
	ID        []byte          `reform:"id,pk"`
	AccountID int64           `reform:"account_id"`
	Settings  json.RawMessage `reform:"settings"`
}

type EndpointConfigTable struct {
	db *sql.DB
}

type FetchReq struct {
	ID []byte `reform:"id,pk"`
}

const fetchQuery = `
SELECT
  id,
  account_id,
  settings
FROM endpoint_config
WHERE
  id = $1
`

func (t *EndpointConfigTable) Fetch(ctx context.Context, req *FetchReq) (*EndpointConfig, error) {
	var ec EndpointConfig
	row := t.db.QueryRowContext(ctx, fetchQuery, req.ID)
	return &ec, row.Scan(&ec.ID, &ec.AccountID, &ec.Settings)
}

type FilterReq struct {
	ID   []byte `reform:"id"`
	Acct int64  `reform:"account_id"`
}

const filterQuery = `
SELECT
  id,
  account_id,
  settings
FROM endpoint_config
WHERE id = $1
  AND account_id = $2
`

func (t *EndpointConfigTable) Filter(ctx context.Context, req *FilterReq) ([]*EndpointConfig, error) {
	var ec EndpointConfig
	row := t.db.QueryRowContext(ctx, fetchQuery, req.ID, req.Acct)
	return []*EndpointConfig{&ec}, row.Scan(&ec.ID, &ec.AccountID, &ec.Settings)
}

func main() {
}
