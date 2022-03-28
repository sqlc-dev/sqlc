// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.12.0

package querytest

import (
	"fmt"
)

// this is the mood type
type FooMood string

const (
	FooMoodSad   FooMood = "sad"
	FooMoodOk    FooMood = "ok"
	FooMoodHappy FooMood = "happy"
)

func (e *FooMood) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FooMood(s)
	case string:
		*e = FooMood(s)
	default:
		return fmt.Errorf("unsupported scan type for FooMood: %T", src)
	}
	return nil
}

// this is the bar table
type FooBar struct {
	// this is the baz column
	Baz string
}
