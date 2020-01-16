// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dinosql

import (
	"fmt"
	"reflect"

	nodes "github.com/lfittl/pg_query_go/nodes"
)

// An ApplyFunc is invoked by Apply for each node n, even if n is nil,
// before and/or after the node's children, using a Cursor describing
// the current node and providing operations on it.
//
// The return value of ApplyFunc controls the syntax tree traversal.
// See Apply for details.
type ApplyFunc func(*Cursor) bool

// Apply traverses a syntax tree recursively, starting with root,
// and calling pre and post for each node as described below.
// Apply returns the syntax tree, possibly modified.
//
// If pre is not nil, it is called for each node before the node's
// children are traversed (pre-order). If pre returns false, no
// children are traversed, and post is not called for that node.
//
// If post is not nil, and a prior call of pre didn't return false,
// post is called for each node after its children are traversed
// (post-order). If post returns false, traversal is terminated and
// Apply returns immediately.
//
// Only fields that refer to AST nodes are considered children;
// i.e., token.Pos, Scopes, Objects, and fields of basic types
// (strings, etc.) are ignored.
//
// Children are traversed in the order in which they appear in the
// respective node's struct definition. A package's files are
// traversed in the filenames' alphabetical order.
//
func Apply(root nodes.Node, pre, post ApplyFunc) (result nodes.Node) {
	parent := &struct{ nodes.Node }{root}
	a := &application{pre: pre, post: post}
	return a.apply(parent, "Node", nil, root)
}

var abort = new(int) // singleton, to signal termination of Apply

// A Cursor describes a node encountered during Apply.
// Information about the node and its parent is available
// from the Node, Parent, Name, and Index methods.
//
// If p is a variable of type and value of the current parent node
// c.Parent(), and f is the field identifier with name c.Name(),
// the following invariants hold:
//
//   p.f            == c.Node()  if c.Index() <  0
//   p.f[c.Index()] == c.Node()  if c.Index() >= 0
//
// The methods Replace, Delete, InsertBefore, and InsertAfter
// can be used to change the AST without disrupting Apply.
type Cursor struct {
	parent nodes.Node
	name   string
	iter   *iterator // valid if non-nil
	node   nodes.Node
}

// Node returns the current Node.
func (c *Cursor) Node() nodes.Node { return c.node }

// Parent returns the parent of the current Node.
func (c *Cursor) Parent() nodes.Node { return c.parent }

// Name returns the name of the parent Node field that contains the current Node.
// If the parent is a *ast.Package and the current Node is a *ast.File, Name returns
// the filename for the current Node.
func (c *Cursor) Name() string { return c.name }

// Index reports the index >= 0 of the current Node in the slice of Nodes that
// contains it, or a value < 0 if the current Node is not part of a slice.
// The index of the current node changes if InsertBefore is called while
// processing the current node.
func (c *Cursor) Index() int {
	if c.iter != nil {
		return c.iter.index
	}
	return -1
}

// field returns the current node's parent field value.
func (c *Cursor) field() reflect.Value {
	return reflect.Indirect(reflect.ValueOf(c.parent)).FieldByName(c.name)
}

// Replace replaces the current Node with n.
// The replacement node is not walked by Apply.
func (c *Cursor) Replace(n nodes.Node) {
	c.node = n
}

// application carries all the shared data so we can pass it around cheaply.
type application struct {
	pre, post ApplyFunc
	cursor    Cursor
	iter      iterator
}

func (a *application) apply(parent nodes.Node, name string, iter *iterator, n nodes.Node) nodes.Node {
	// avoid heap-allocating a new cursor for each apply call; reuse a.cursor instead
	cursor := Cursor{
		parent: parent,
		name:   name,
		iter:   iter,
		node:   n,
	}

	if a.pre != nil && !a.pre(&cursor) {
		return cursor.node
	}

	// walk children
	// (the order of the cases matches the order of the corresponding node types in go/ast)
	switch n := n.(type) {
	case nil:
		// nothing to do

	case nodes.RawStmt:
		n.Stmt = a.apply(n, "Stmt", nil, n.Stmt)
		cursor.node = n

	case nodes.SelectStmt:
		n.TargetList = a.apply(n, "TargetList", nil, n.TargetList)
		cursor.node = n

	default:
		panic(fmt.Sprintf("Apply: unexpected node type %T", n))
	}

	if a.post != nil && !a.post(&cursor) {
		panic(abort)
	}

	return cursor.node
}

// An iterator controls iteration over a slice of nodes.
type iterator struct {
	index, step int
}

func (a *application) applyList(parent nodes.Node, name string) {
	// avoid heap-allocating a new iterator for each applyList call; reuse a.iter instead
	saved := a.iter
	a.iter.index = 0
	for {
		// must reload parent.name each time, since cursor modifications might change it
		v := reflect.Indirect(reflect.ValueOf(parent)).FieldByName(name)
		if a.iter.index >= v.Len() {
			break
		}

		// element x may be nil in a bad AST - be cautious
		var x nodes.Node
		if e := v.Index(a.iter.index); e.IsValid() {
			x = e.Interface().(nodes.Node)
		}

		a.iter.step = 1
		a.apply(parent, name, &a.iter, x)
		a.iter.index += a.iter.step
	}
	a.iter = saved
}
