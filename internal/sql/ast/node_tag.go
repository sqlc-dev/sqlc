package ast

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// NodeTag names its containing node type in JSON.
//
// Node is an interface, so the JSON encoding of an AST would otherwise carry no
// record of which node a given object is. Nodes with no fields of their own
// (A_Star, Null, TODO) would all encode as "{}", and an empty List would be
// indistinguishable from them. Every node struct instead declares
//
//	Tag NodeTag[T] `json:"tag"`
//
// as its first field, with T the node's own type. The field is zero-sized and
// its zero value is valid, so constructing nodes as struct literals is
// unaffected; encoding/json calls MarshalJSON on the field, and the type
// parameter carries the node's identity to it at compile time. A test checks
// that every node declares the field and that its type parameter matches.
type NodeTag[T any] struct{}

func (NodeTag[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(reflect.TypeFor[T]().Name())
}

// UnmarshalJSON accepts only the containing node's own name, so a document
// decoded into the wrong node type is an error rather than a silently
// misfielded tree.
func (*NodeTag[T]) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	if want := reflect.TypeFor[T]().Name(); name != want {
		return fmt.Errorf("node tagged %q decoded into %s", name, want)
	}
	return nil
}
