package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// TagKey is the JSON key under which a node's concrete type is reported.
//
// Node is an interface, so the JSON encoding of an AST would otherwise carry no
// record of which node a given object is. Nodes with no fields (A_Star, Null,
// TODO) all encode as "{}", and an empty List is indistinguishable from them.
// Every node object is emitted with this key first, holding the name of its Go
// type.
const TagKey = "tag"

var nodeType = reflect.TypeOf((*Node)(nil)).Elem()

// MarshalJSON encodes the statement and every node beneath it, tagging each
// node object with its type. RawStmt is the root of the AST that the parse and
// analyze commands print, so implementing it here is enough to tag a whole
// tree.
func (n *RawStmt) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	return marshalValue(reflect.ValueOf(n))
}

func marshalValue(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return []byte("null"), nil

	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			return []byte("null"), nil
		}
		return marshalValue(v.Elem())

	case reflect.Struct:
		return marshalStruct(v)

	case reflect.Slice:
		if v.IsNil() {
			return []byte("null"), nil
		}
		fallthrough
	case reflect.Array:
		return marshalArray(v)

	default:
		// Scalars, strings and anything else encoding/json already handles.
		return json.Marshal(v.Interface())
	}
}

func marshalStruct(v reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	// A node's Pos method is declared on the pointer type.
	if reflect.PointerTo(v.Type()).Implements(nodeType) {
		fmt.Fprintf(&buf, "%q:%q", TagKey, v.Type().Name())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			tagName, _, _ := strings.Cut(tag, ",")
			if tagName == "-" {
				continue
			}
			if tagName != "" {
				name = tagName
			}
		}
		value, err := marshalValue(v.Field(i))
		if err != nil {
			return nil, err
		}
		if buf.Len() > 1 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(&buf, "%q:", name)
		buf.Write(value)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func marshalArray(v reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		item, err := marshalValue(v.Index(i))
		if err != nil {
			return nil, err
		}
		buf.Write(item)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}
