// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"database/sql/driver"
	"fmt"
)

type QueryParamEnumTableEnum string

const (
	QueryParamEnumTableEnumG QueryParamEnumTableEnum = "g"
	QueryParamEnumTableEnumH QueryParamEnumTableEnum = "h"
)

func (e *QueryParamEnumTableEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QueryParamEnumTableEnum(s)
	case string:
		*e = QueryParamEnumTableEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for QueryParamEnumTableEnum: %T", src)
	}
	return nil
}

type NullQueryParamEnumTableEnum struct {
	QueryParamEnumTableEnum QueryParamEnumTableEnum
	Valid                   bool // Valid is true if QueryParamEnumTableEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQueryParamEnumTableEnum) Scan(value interface{}) error {
	if value == nil {
		ns.QueryParamEnumTableEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QueryParamEnumTableEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQueryParamEnumTableEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QueryParamEnumTableEnum), nil
}

type QueryParamStructEnumTableEnum string

const (
	QueryParamStructEnumTableEnumI QueryParamStructEnumTableEnum = "i"
	QueryParamStructEnumTableEnumJ QueryParamStructEnumTableEnum = "j"
)

func (e *QueryParamStructEnumTableEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QueryParamStructEnumTableEnum(s)
	case string:
		*e = QueryParamStructEnumTableEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for QueryParamStructEnumTableEnum: %T", src)
	}
	return nil
}

type NullQueryParamStructEnumTableEnum struct {
	QueryParamStructEnumTableEnum QueryParamStructEnumTableEnum
	Valid                         bool // Valid is true if QueryParamStructEnumTableEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQueryParamStructEnumTableEnum) Scan(value interface{}) error {
	if value == nil {
		ns.QueryParamStructEnumTableEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QueryParamStructEnumTableEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQueryParamStructEnumTableEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QueryParamStructEnumTableEnum), nil
}

type QueryReturnEnumTableEnum string

const (
	QueryReturnEnumTableEnumK QueryReturnEnumTableEnum = "k"
	QueryReturnEnumTableEnumL QueryReturnEnumTableEnum = "l"
)

func (e *QueryReturnEnumTableEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QueryReturnEnumTableEnum(s)
	case string:
		*e = QueryReturnEnumTableEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for QueryReturnEnumTableEnum: %T", src)
	}
	return nil
}

type NullQueryReturnEnumTableEnum struct {
	QueryReturnEnumTableEnum QueryReturnEnumTableEnum
	Valid                    bool // Valid is true if QueryReturnEnumTableEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQueryReturnEnumTableEnum) Scan(value interface{}) error {
	if value == nil {
		ns.QueryReturnEnumTableEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QueryReturnEnumTableEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQueryReturnEnumTableEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QueryReturnEnumTableEnum), nil
}

type QueryReturnFullTableEnum string

const (
	QueryReturnFullTableEnumE QueryReturnFullTableEnum = "e"
	QueryReturnFullTableEnumF QueryReturnFullTableEnum = "f"
)

func (e *QueryReturnFullTableEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QueryReturnFullTableEnum(s)
	case string:
		*e = QueryReturnFullTableEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for QueryReturnFullTableEnum: %T", src)
	}
	return nil
}

type NullQueryReturnFullTableEnum struct {
	QueryReturnFullTableEnum QueryReturnFullTableEnum
	Valid                    bool // Valid is true if QueryReturnFullTableEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQueryReturnFullTableEnum) Scan(value interface{}) error {
	if value == nil {
		ns.QueryReturnFullTableEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QueryReturnFullTableEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQueryReturnFullTableEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QueryReturnFullTableEnum), nil
}

type QueryReturnStructEnumTableEnum string

const (
	QueryReturnStructEnumTableEnumK QueryReturnStructEnumTableEnum = "k"
	QueryReturnStructEnumTableEnumL QueryReturnStructEnumTableEnum = "l"
)

func (e *QueryReturnStructEnumTableEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QueryReturnStructEnumTableEnum(s)
	case string:
		*e = QueryReturnStructEnumTableEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for QueryReturnStructEnumTableEnum: %T", src)
	}
	return nil
}

type NullQueryReturnStructEnumTableEnum struct {
	QueryReturnStructEnumTableEnum QueryReturnStructEnumTableEnum
	Valid                          bool // Valid is true if QueryReturnStructEnumTableEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQueryReturnStructEnumTableEnum) Scan(value interface{}) error {
	if value == nil {
		ns.QueryReturnStructEnumTableEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QueryReturnStructEnumTableEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQueryReturnStructEnumTableEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QueryReturnStructEnumTableEnum), nil
}

type QueryParamEnumTable struct {
	ID    int32
	Other QueryParamEnumTableEnum
	Value NullQueryParamEnumTableEnum
}

type QueryReturnFullTable struct {
	ID    int32
	Value NullQueryReturnFullTableEnum
}
