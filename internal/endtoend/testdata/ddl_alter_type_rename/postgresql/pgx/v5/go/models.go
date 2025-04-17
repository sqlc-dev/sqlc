// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"database/sql/driver"
	"fmt"
)

type NewEvent string

const (
	NewEventSTART NewEvent = "START"
	NewEventSTOP  NewEvent = "STOP"
)

func (e *NewEvent) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NewEvent(s)
	case string:
		*e = NewEvent(s)
	default:
		return fmt.Errorf("unsupported scan type for NewEvent: %T", src)
	}
	return nil
}

type NullNewEvent struct {
	NewEvent NewEvent
	Valid    bool // Valid is true if NewEvent is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullNewEvent) Scan(value interface{}) error {
	if value == nil {
		ns.NewEvent, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.NewEvent.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullNewEvent) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.NewEvent), nil
}

type LogLine struct {
	ID     int64
	Status NewEvent
}
