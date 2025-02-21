// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package querytest

import (
	"context"
)

const getAllOrganisations = `-- name: GetAllOrganisations :many
SELECT party_id, name, legal_name FROM organisation
`

func (q *Queries) GetAllOrganisations(ctx context.Context) ([]Organisation, error) {
	rows, err := q.db.QueryContext(ctx, getAllOrganisations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Organisation
	for rows.Next() {
		var i Organisation
		if err := rows.Scan(&i.PartyID, &i.Name, &i.LegalName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllParties = `-- name: GetAllParties :many
SELECT party_id, name FROM party
`

func (q *Queries) GetAllParties(ctx context.Context) ([]Party, error) {
	rows, err := q.db.QueryContext(ctx, getAllParties)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Party
	for rows.Next() {
		var i Party
		if err := rows.Scan(&i.PartyID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllPeople = `-- name: GetAllPeople :many
SELECT party_id, name, first_name, last_name FROM person
`

func (q *Queries) GetAllPeople(ctx context.Context) ([]Person, error) {
	rows, err := q.db.QueryContext(ctx, getAllPeople)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Person
	for rows.Next() {
		var i Person
		if err := rows.Scan(
			&i.PartyID,
			&i.Name,
			&i.FirstName,
			&i.LastName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
