// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
	"database/sql"

	"github.com/kyleconroy/sqlc-testdata/mysql"
)

const getAll = `-- name: GetAll :many
select id, first_name, last_name, age, job_status, created from users
`

func (q *Queries) GetAll(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Age,
			&i.JobStatus,
			&i.Created,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUsersOrders = `-- name: GetAllUsersOrders :many
select u.id as user_id, u.first_name, o.price, o.id as order_id from orders as o left join users as u on u.id = o.user_id
`

type GetAllUsersOrdersRow struct {
	UserID    sql.NullInt64
	FirstName sql.NullString
	Price     float64
	OrderID   int
}

func (q *Queries) GetAllUsersOrders(ctx context.Context) ([]GetAllUsersOrdersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsersOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersOrdersRow
	for rows.Next() {
		var i GetAllUsersOrdersRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.Price,
			&i.OrderID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCount = `-- name: GetCount :one
select id as my_id, COUNT(id) as id_count from users where id > 4
`

type GetCountRow struct {
	MyID    int
	IDCount int
}

func (q *Queries) GetCount(ctx context.Context) (GetCountRow, error) {
	row := q.db.QueryRowContext(ctx, getCount)
	var i GetCountRow
	err := row.Scan(&i.MyID, &i.IDCount)
	return i, err
}

const getNameByID = `-- name: GetNameByID :one
select first_name, last_name from users where id = ?
`

type GetNameByIDRow struct {
	FirstName string
	LastName  sql.NullString
}

func (q *Queries) GetNameByID(ctx context.Context, id mysql.ID) (GetNameByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getNameByID, id)
	var i GetNameByIDRow
	err := row.Scan(&i.FirstName, &i.LastName)
	return i, err
}

const insertNewUser = `-- name: InsertNewUser :exec
insert into users(first_name, last_name) values (?, ?)
`

type InsertNewUserParams struct {
	FirstName string
	LastName  sql.NullString
}

func (q *Queries) InsertNewUser(ctx context.Context, arg InsertNewUserParams) error {
	_, err := q.db.ExecContext(ctx, insertNewUser, arg.FirstName, arg.LastName)
	return err
}

const insertUsersFromOrders = `-- name: InsertUsersFromOrders :exec
insert into users(first_name) select user_id from orders where id = ?
`

func (q *Queries) InsertUsersFromOrders(ctx context.Context, id mysql.ID) error {
	_, err := q.db.ExecContext(ctx, insertUsersFromOrders, id)
	return err
}

const updateAllUsers = `-- name: UpdateAllUsers :exec
update users set first_name = 'Bob'
`

func (q *Queries) UpdateAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, updateAllUsers)
	return err
}

const updateUserAt = `-- name: UpdateUserAt :exec
update users set first_name = ?, last_name = ? where id > ? and first_name = ? limit 3
`

type UpdateUserAtParams struct {
	FirstName   string
	LastName    sql.NullString
	ID          mysql.ID
	FirstName_2 string
}

func (q *Queries) UpdateUserAt(ctx context.Context, arg UpdateUserAtParams) error {
	_, err := q.db.ExecContext(ctx, updateUserAt,
		arg.FirstName,
		arg.LastName,
		arg.ID,
		arg.FirstName_2,
	)
	return err
}
