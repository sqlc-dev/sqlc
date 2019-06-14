package equinox

import (
	"context"
	"database/sql"
	"time"
)

type QueryRow interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// CREATE TABLE account (
//     id SERIAL UNIQUE NOT NULL,
//     created_at timestamp without time zone NOT NULL,
//     name character varying(255) NOT NULL,
//     stripe_id character varying(64),
//     plan character varying(64),
//     slug character varying(255) NOT NULL,
//     owner_id integer NOT NULL
// );

type Account struct {
	ID        int
	CreatedAt time.Time
	Name      string
	StripeID  sql.NullString
	Plan      sql.NullString
	Slug      string
	OwnerID   int
}

const getAccountBySlug = `
SELECT 
  id,
  created_at,
  name,
  stripe_id, 
  plan,
  slug,
  owner_id
FROM account
WHERE slug = $1
`

func GetAccountBySlug(q QueryRow, ctx context.Context, slug string) (Account, error) {
	var a Account
	row := q.QueryRowContext(ctx, getAccountBySlug, slug)
	return a, row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.Name,
		&a.StripeID,
		&a.Plan,
		&a.Slug,
		&a.OwnerID,
	)
}

const getAccountByID = `
SELECT 
  id,
  created_at,
  name,
  stripe_id, 
  plan,
  slug,
  owner_id
FROM account
WHERE id = $1
`

func GetAccountByID(q QueryRow, ctx context.Context, id int) (Account, error) {
	var a Account
	row := q.QueryRowContext(ctx, getAccountByID, id)
	return a, scanAccount(row, &a)
}

const getAccountByUser = `
SELECT 
  id,
  created_at,
  name,
  stripe_id, 
  plan,
  slug,
  owner_id
FROM account
WHERE slug = $1
AND EXISTS (
	SELECT user_id
	FROM membership
	WHERE account_id = account.id AND user_id = $2
)`

func GetAccountByUser(q QueryRow, ctx context.Context, slug string, userID int) (Account, error) {
	var a Account
	row := q.QueryRowContext(ctx, getAccountByUser, slug, userID)
	return a, scanAccount(row, &a)
}

const getDefaultAccountForUser = `
SELECT 
  id,
  created_at,
  name,
  stripe_id, 
  plan,
  slug,
  owner_id
FROM account
WHERE id IN (
    SELECT account_id
	FROM "user", membership
	WHERE "user".default_membership_id = membership.id
    AND "user".id = $1
)`

func GetDefaultAccountForUser(q QueryRow, ctx context.Context, userID int) (Account, error) {
	var a Account
	row := q.QueryRowContext(ctx, getDefaultAccountForUser, userID)
	return a, scanAccount(row, &a)
}
