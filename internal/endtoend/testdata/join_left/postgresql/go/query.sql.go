// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const allAuthors = `-- name: AllAuthors :many
SELECT  a.id, a.name, a.parent_id, p.id, p.name, p.parent_id
FROM    authors a
        LEFT JOIN authors p
            ON a.parent_id = p.id
`

type AllAuthorsRow struct {
	ID         int32
	Name       string
	ParentID   sql.NullInt32
	ID_2       sql.NullInt32
	Name_2     sql.NullString
	ParentID_2 sql.NullInt32
}

func (q *Queries) AllAuthors(ctx context.Context) ([]AllAuthorsRow, error) {
	ctx, done := q.observer(ctx, "AllAuthors")
	rows, err := q.db.QueryContext(ctx, allAuthors)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []AllAuthorsRow
	for rows.Next() {
		var i AllAuthorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.ID_2,
			&i.Name_2,
			&i.ParentID_2,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const allAuthorsAliases = `-- name: AllAuthorsAliases :many
SELECT  a.id, a.name, a.parent_id, p.id, p.name, p.parent_id
FROM    authors a
        LEFT JOIN authors p
            ON a.parent_id = p.id
`

type AllAuthorsAliasesRow struct {
	ID         int32
	Name       string
	ParentID   sql.NullInt32
	ID_2       sql.NullInt32
	Name_2     sql.NullString
	ParentID_2 sql.NullInt32
}

func (q *Queries) AllAuthorsAliases(ctx context.Context) ([]AllAuthorsAliasesRow, error) {
	ctx, done := q.observer(ctx, "AllAuthorsAliases")
	rows, err := q.db.QueryContext(ctx, allAuthorsAliases)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []AllAuthorsAliasesRow
	for rows.Next() {
		var i AllAuthorsAliasesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.ID_2,
			&i.Name_2,
			&i.ParentID_2,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const allAuthorsAliases2 = `-- name: AllAuthorsAliases2 :many
SELECT  a.id, a.name, a.parent_id, p.id, p.name, p.parent_id
FROM    authors a
        LEFT JOIN authors p
            ON a.parent_id = p.id
`

type AllAuthorsAliases2Row struct {
	ID         int32
	Name       string
	ParentID   sql.NullInt32
	ID_2       sql.NullInt32
	Name_2     sql.NullString
	ParentID_2 sql.NullInt32
}

func (q *Queries) AllAuthorsAliases2(ctx context.Context) ([]AllAuthorsAliases2Row, error) {
	ctx, done := q.observer(ctx, "AllAuthorsAliases2")
	rows, err := q.db.QueryContext(ctx, allAuthorsAliases2)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []AllAuthorsAliases2Row
	for rows.Next() {
		var i AllAuthorsAliases2Row
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.ID_2,
			&i.Name_2,
			&i.ParentID_2,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const allSuperAuthors = `-- name: AllSuperAuthors :many
SELECT  id, name, parent_id, super_id, super_name, super_parent_id
FROM    authors
        LEFT JOIN super_authors
            ON authors.parent_id = super_authors.super_id
`

type AllSuperAuthorsRow struct {
	ID            int32
	Name          string
	ParentID      sql.NullInt32
	SuperID       sql.NullInt32
	SuperName     sql.NullString
	SuperParentID sql.NullInt32
}

func (q *Queries) AllSuperAuthors(ctx context.Context) ([]AllSuperAuthorsRow, error) {
	ctx, done := q.observer(ctx, "AllSuperAuthors")
	rows, err := q.db.QueryContext(ctx, allSuperAuthors)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []AllSuperAuthorsRow
	for rows.Next() {
		var i AllSuperAuthorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.SuperID,
			&i.SuperName,
			&i.SuperParentID,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const allSuperAuthorsAliases = `-- name: AllSuperAuthorsAliases :many
SELECT  id, name, parent_id, super_id, super_name, super_parent_id
FROM    authors a
        LEFT JOIN super_authors sa
            ON a.parent_id = sa.super_id
`

type AllSuperAuthorsAliasesRow struct {
	ID            int32
	Name          string
	ParentID      sql.NullInt32
	SuperID       sql.NullInt32
	SuperName     sql.NullString
	SuperParentID sql.NullInt32
}

func (q *Queries) AllSuperAuthorsAliases(ctx context.Context) ([]AllSuperAuthorsAliasesRow, error) {
	ctx, done := q.observer(ctx, "AllSuperAuthorsAliases")
	rows, err := q.db.QueryContext(ctx, allSuperAuthorsAliases)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []AllSuperAuthorsAliasesRow
	for rows.Next() {
		var i AllSuperAuthorsAliasesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.SuperID,
			&i.SuperName,
			&i.SuperParentID,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const allSuperAuthorsAliases2 = `-- name: AllSuperAuthorsAliases2 :many
SELECT  a.id, a.name, a.parent_id, sa.super_id, sa.super_name, sa.super_parent_id
FROM    authors a
        LEFT JOIN super_authors sa
            ON a.parent_id = sa.super_id
`

type AllSuperAuthorsAliases2Row struct {
	ID            int32
	Name          string
	ParentID      sql.NullInt32
	SuperID       sql.NullInt32
	SuperName     sql.NullString
	SuperParentID sql.NullInt32
}

func (q *Queries) AllSuperAuthorsAliases2(ctx context.Context) ([]AllSuperAuthorsAliases2Row, error) {
	ctx, done := q.observer(ctx, "AllSuperAuthorsAliases2")
	rows, err := q.db.QueryContext(ctx, allSuperAuthorsAliases2)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []AllSuperAuthorsAliases2Row
	for rows.Next() {
		var i AllSuperAuthorsAliases2Row
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.SuperID,
			&i.SuperName,
			&i.SuperParentID,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const getMayors = `-- name: GetMayors :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
INNER JOIN mayors USING (mayor_id)
`

type GetMayorsRow struct {
	UserID   int32
	FullName string
}

func (q *Queries) GetMayors(ctx context.Context) ([]GetMayorsRow, error) {
	ctx, done := q.observer(ctx, "GetMayors")
	rows, err := q.db.QueryContext(ctx, getMayors)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []GetMayorsRow
	for rows.Next() {
		var i GetMayorsRow
		if err := rows.Scan(&i.UserID, &i.FullName); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const getMayorsOptional = `-- name: GetMayorsOptional :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
LEFT JOIN mayors USING (mayor_id)
`

type GetMayorsOptionalRow struct {
	UserID   int32
	FullName sql.NullString
}

func (q *Queries) GetMayorsOptional(ctx context.Context) ([]GetMayorsOptionalRow, error) {
	ctx, done := q.observer(ctx, "GetMayorsOptional")
	rows, err := q.db.QueryContext(ctx, getMayorsOptional)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []GetMayorsOptionalRow
	for rows.Next() {
		var i GetMayorsOptionalRow
		if err := rows.Scan(&i.UserID, &i.FullName); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}

const getSuggestedUsersByID = `-- name: GetSuggestedUsersByID :many
SELECT  DISTINCT u.user_id, u.user_nickname, u.user_email, u.user_display_name, u.user_password, u.user_google_id, u.user_apple_id, u.user_bio, u.user_created_at, u.user_avatar_id, m.media_id, m.media_created_at, m.media_hash, m.media_directory, m.media_author_id, m.media_width, m.media_height
FROM    users_2 u
        LEFT JOIN media m
            ON u.user_avatar_id = m.media_id
WHERE   u.user_id != $1
LIMIT   $2
`

type GetSuggestedUsersByIDParams struct {
	UserID   uuid.UUID
	UserImit int32
}

type GetSuggestedUsersByIDRow struct {
	UserID          uuid.UUID
	UserNickname    string
	UserEmail       string
	UserDisplayName string
	UserPassword    sql.NullString
	UserGoogleID    sql.NullString
	UserAppleID     sql.NullString
	UserBio         string
	UserCreatedAt   time.Time
	UserAvatarID    uuid.NullUUID
	MediaID         uuid.NullUUID
	MediaCreatedAt  sql.NullTime
	MediaHash       sql.NullString
	MediaDirectory  sql.NullString
	MediaAuthorID   uuid.NullUUID
	MediaWidth      sql.NullInt32
	MediaHeight     sql.NullInt32
}

func (q *Queries) GetSuggestedUsersByID(ctx context.Context, arg GetSuggestedUsersByIDParams) ([]GetSuggestedUsersByIDRow, error) {
	ctx, done := q.observer(ctx, "GetSuggestedUsersByID")
	rows, err := q.db.QueryContext(ctx, getSuggestedUsersByID, arg.UserID, arg.UserImit)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []GetSuggestedUsersByIDRow
	for rows.Next() {
		var i GetSuggestedUsersByIDRow
		if err := rows.Scan(
			&i.UserID,
			&i.UserNickname,
			&i.UserEmail,
			&i.UserDisplayName,
			&i.UserPassword,
			&i.UserGoogleID,
			&i.UserAppleID,
			&i.UserBio,
			&i.UserCreatedAt,
			&i.UserAvatarID,
			&i.MediaID,
			&i.MediaCreatedAt,
			&i.MediaHash,
			&i.MediaDirectory,
			&i.MediaAuthorID,
			&i.MediaWidth,
			&i.MediaHeight,
		); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
