package equinox

import (
	"github.com/kyleconroy/strongdb/examples/equinox/credentials"
)

// func (c *Client) initCred() {
// 	const SELECT = `
// 	SELECT id, sid, created, accountid, tokenhash
// 	FROM credentials`
//
// 	c.q.listCreds = c.prepareNamed(SELECT + `
// 	WHERE
// 		id > :after AND accountid = :accountid
// 	ORDER BY id desc
// 	LIMIT :limit
// 	`)
//

// 	SELECT id, sid, created, accountid, tokenhash
// 	FROM credentials
// 	WHERE sid = $1
var getCredBySID = credentials.Select{
	Where: credentials.Eq{credentials.SID, credentials.Arg},
}

//
// 	c.q.createCred = c.prepare(`
// 	INSERT INTO credentials
// 	(accountid, tokenhash, sid) VALUES ($1, $2, $3)
// 	RETURNING id, sid, created, accountid, tokenhash
// 	`)
//
// 	c.q.scopedDeleteCredBySID = c.prepare(`
// 	DELETE FROM credentials
// 	WHERE
// 		accountid = $1 AND sid = $2
// 	`)
// }
