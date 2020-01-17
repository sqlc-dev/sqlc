// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

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
	defer func() {
		if r := recover(); r != nil && r != abort {
			panic(r)
		}
		result = parent.Node
	}()
	a := &application{pre: pre, post: post}
	a.apply(parent, "Node", nil, root)
	return
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
	v := c.field()
	if i := c.Index(); i >= 0 {
		v = v.Index(i)
	}
	v.Set(reflect.ValueOf(n))
}

// application carries all the shared data so we can pass it around cheaply.
type application struct {
	pre, post ApplyFunc
	cursor    Cursor
	iter      iterator
}

func (a *application) apply(parent nodes.Node, name string, iter *iterator, node nodes.Node) {
	// convert typed nil into untyped nil
	if v := reflect.ValueOf(node); v.Kind() == reflect.Ptr && v.IsNil() {
		node = nil
	}

	// avoid heap-allocating a new cursor for each apply call; reuse a.cursor instead
	saved := a.cursor
	a.cursor.parent = parent
	a.cursor.name = name
	a.cursor.iter = iter
	a.cursor.node = node

	if a.pre != nil && !a.pre(&a.cursor) {
		a.cursor = saved
		return
	}

	// walk children
	// (the order of the cases matches the order of the corresponding node types in go/ast)
	switch n := node.(type) {
	case nil:
		// nothing to do

	case nodes.A_ArrayExpr:
		a.apply(&n, "Elements", nil, n.Elements)
		a.cursor.Replace(n)

	case nodes.A_Const:
		a.apply(&n, "Val", nil, n.Val)
		a.cursor.Replace(n)

	case nodes.A_Expr:
		a.apply(&n, "Name", nil, n.Name)
		a.apply(&n, "Lexpr", nil, n.Lexpr)
		a.apply(&n, "Rexpr", nil, n.Rexpr)
		a.cursor.Replace(n)

	case nodes.A_Indices:
		a.apply(&n, "Lidx", nil, n.Lidx)
		a.apply(&n, "Uidx", nil, n.Uidx)
		a.cursor.Replace(n)

	case nodes.A_Indirection:
		a.apply(&n, "Arg", nil, n.Arg)
		a.apply(&n, "Indirection", nil, n.Indirection)
		a.cursor.Replace(n)

	case nodes.A_Star:
		// pass

	case nodes.AccessPriv:
		a.apply(&n, "Cols", nil, n.Cols)
		a.cursor.Replace(n)

	case nodes.Aggref:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Aggargtypes", nil, n.Aggargtypes)
		a.apply(&n, "Aggdirectargs", nil, n.Aggdirectargs)
		a.apply(&n, "Args", nil, n.Args)
		a.apply(&n, "Aggorder", nil, n.Aggorder)
		a.apply(&n, "Aggdistinct", nil, n.Aggdistinct)
		a.apply(&n, "Aggfilter", nil, n.Aggfilter)
		a.cursor.Replace(n)

	case nodes.Alias:
		a.apply(&n, "Colnames", nil, n.Colnames)
		a.cursor.Replace(n)

	case nodes.AlterCollationStmt:
		a.apply(&n, "Collname", nil, n.Collname)
		a.cursor.Replace(n)

	case nodes.AlterDatabaseSetStmt:
		if n.Setstmt != nil {
			a.apply(&n, "Setstmt", nil, *n.Setstmt)
			a.cursor.Replace(n)
		}

	case nodes.AlterDatabaseStmt:
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterDefaultPrivilegesStmt:
		if n.Action != nil {
			a.apply(&n, "Action", nil, *n.Action)
		}
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterDomainStmt:
		a.apply(&n, "TypeName", nil, n.TypeName)
		a.apply(&n, "Def", nil, n.Def)
		a.cursor.Replace(n)

	case nodes.AlterEnumStmt:
		a.apply(&n, "TypeName", nil, n.TypeName)
		a.cursor.Replace(n)

	case nodes.AlterEventTrigStmt:
		// pass

	case nodes.AlterExtensionContentsStmt:
		a.apply(&n, "Object", nil, n.Object)
		a.cursor.Replace(n)

	case nodes.AlterExtensionStmt:
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterFdwStmt:
		a.apply(&n, "FuncOptions", nil, n.FuncOptions)
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterForeignServerStmt:
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterFunctionStmt:
		if n.Func != nil {
			a.apply(&n, "Func", nil, n.Func)
		}
		a.apply(&n, "Actions", nil, n.Actions)
		a.cursor.Replace(n)

	case nodes.AlterObjectDependsStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Object", nil, n.Object)
		a.apply(&n, "Extname", nil, n.Extname)
		a.cursor.Replace(n)

	case nodes.AlterObjectSchemaStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Object", nil, n.Object)
		a.cursor.Replace(n)

	case nodes.AlterOpFamilyStmt:
		a.apply(&n, "Opfamilyname", nil, n.Opfamilyname)
		a.apply(&n, "Items", nil, n.Items)
		a.cursor.Replace(n)

	case nodes.AlterOperatorStmt:
		if n.Opername != nil {
			a.apply(&n, "Opername", nil, *n.Opername)
		}
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterOwnerStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Object", nil, n.Object)
		if n.Newowner != nil {
			a.apply(&n, "Newowner", nil, *n.Newowner)
		}
		a.cursor.Replace(n)

	case nodes.AlterPolicyStmt:
		if n.Table != nil {
			a.apply(&n, "Table", nil, *n.Table)
		}
		a.apply(&n, "Roles", nil, n.Roles)
		a.apply(&n, "Qual", nil, n.Qual)
		a.apply(&n, "WithCheck", nil, n.WithCheck)
		a.cursor.Replace(n)

	case nodes.AlterPublicationStmt:
		a.apply(&n, "Options", nil, n.Options)
		a.apply(&n, "Table", nil, n.Tables)
		a.cursor.Replace(n)

	case nodes.AlterRoleSetStmt:
		if n.Role != nil {
			a.apply(&n, "Role", nil, *n.Role)
		}
		a.apply(&n, "Setstmt", nil, n.Setstmt)
		a.cursor.Replace(n)

	case nodes.AlterRoleStmt:
		if n.Role != nil {
			a.apply(&n, "Role", nil, *n.Role)
		}
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterSeqStmt:
		if n.Sequence != nil {
			a.apply(&n, "Sequence", nil, *n.Sequence)
		}
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterSubscriptionStmt:
		a.apply(&n, "Publication", nil, n.Publication)
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterSystemStmt:
		a.apply(&n, "Setstmt", nil, n.Setstmt)
		a.cursor.Replace(n)

	case nodes.AlterTSConfigurationStmt:
		a.apply(&n, "Cfgname", nil, n.Cfgname)
		a.apply(&n, "Tokentype", nil, n.Tokentype)
		a.apply(&n, "Dicts", nil, n.Dicts)
		a.cursor.Replace(n)

	case nodes.AlterTSDictionaryStmt:
		a.apply(&n, "Dictname", nil, n.Dictname)
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterTableCmd:
		if n.Newowner != nil {
			a.apply(&n, "Newowner", nil, *n.Newowner)
		}
		a.apply(&n, "Def", nil, n.Def)
		a.cursor.Replace(n)

	case nodes.AlterTableMoveAllStmt:
		a.apply(&n, "Roles", nil, n.Roles)
		a.cursor.Replace(n)

	case nodes.AlterTableSpaceOptionsStmt:
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlterTableStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Cmds", nil, n.Cmds)
		a.cursor.Replace(n)

	case nodes.AlterUserMappingStmt:
		if n.User != nil {
			a.apply(&n, "User", nil, *n.User)
		}
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.AlternativeSubPlan:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Subplans", nil, n.Subplans)
		a.cursor.Replace(n)

	case nodes.ArrayCoerceExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.ArrayExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Elements", nil, n.Elements)
		a.cursor.Replace(n)

	case nodes.ArrayRef:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Refupperindexpr", nil, n.Refupperindexpr)
		a.apply(&n, "Reflowerindexpr", nil, n.Reflowerindexpr)
		a.apply(&n, "Refexpr", nil, n.Refexpr)
		a.apply(&n, "Refassgnexpr", nil, n.Refassgnexpr)
		a.cursor.Replace(n)

	case nodes.BitString:
		// pass

	case nodes.BlockIdData:
		// pass

	case nodes.BoolExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Args)
		a.cursor.Replace(n)

	case nodes.BooleanTest:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.CaseExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.apply(&n, "Args", nil, n.Args)
		a.apply(&n, "Defresult", nil, n.Defresult)
		a.cursor.Replace(n)

	case nodes.CaseTestExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.cursor.Replace(n)

	case nodes.CaseWhen:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Expr", nil, n.Expr)
		a.apply(&n, "Result", nil, n.Result)
		a.cursor.Replace(n)

	case nodes.CheckPointStmt:
		// pass

	case nodes.ClosePortalStmt:
		// pass

	case nodes.ClusterStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
			a.cursor.Replace(n)
		}

	case nodes.CoalesceExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Args)
		a.cursor.Replace(n)

	case nodes.CoerceToDomain:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.CoerceToDomainValue:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.cursor.Replace(n)

	case nodes.CoerceViaIO:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.CollateClause:
		a.apply(&n, "Arg", nil, n.Arg)
		a.apply(&n, "Collname", nil, n.Collname)
		a.cursor.Replace(n)
	case nodes.CollateExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.ColumnDef:
		if n.TypeName != nil {
			a.apply(&n, "TypeName", nil, *n.TypeName)
		}
		a.apply(&n, "RawDefault", nil, n.RawDefault)
		a.apply(&n, "CookedDefault", nil, n.CookedDefault)
		a.apply(&n, "Constraints", nil, n.Constraints)
		a.apply(&n, "Fdwoptions", nil, n.Fdwoptions)
		a.cursor.Replace(n)

	case nodes.ColumnRef:
		a.apply(&n, "Fields", nil, n.Fields)
		a.cursor.Replace(n)

	case nodes.CommentStmt:
		a.apply(&n, "Object", nil, n.Object)
		a.cursor.Replace(n)

	case nodes.CommonTableExpr:
		a.apply(&n, "Aliascolnames", nil, n.Aliascolnames)
		a.apply(&n, "Ctequery", nil, n.Ctequery)
		a.apply(&n, "Ctecolnames", nil, n.Ctecolnames)
		a.apply(&n, "Ctecolcollations", nil, n.Ctecolcollations)
		a.cursor.Replace(n)

	case nodes.CompositeTypeStmt:
		if n.Typevar != nil {
			a.apply(&n, "Typevar", nil, *n.Typevar)
		}
		a.apply(&n, "Coldeflist", nil, n.Coldeflist)
		a.cursor.Replace(n)

	case nodes.Const:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.cursor.Replace(n)

	case nodes.Constraint:
		a.apply(&n, "RawExpr", nil, n.RawExpr)
		a.apply(&n, "Keys", nil, n.Keys)
		a.apply(&n, "Exclusions", nil, n.Exclusions)
		a.apply(&n, "Options", nil, n.Options)
		a.apply(&n, "WhereClause", nil, n.WhereClause)
		if n.Pktable != nil {
			a.apply(&n, "Pktable", nil, *n.Pktable)
		}
		a.apply(&n, "FkAttrs", nil, n.FkAttrs)
		a.apply(&n, "PkAttrs", nil, n.PkAttrs)
		a.apply(&n, "OldConpfeqop", nil, n.OldConpfeqop)
		a.cursor.Replace(n)

	case nodes.ConstraintsSetStmt:
		a.apply(&n, "Constraints", nil, n.Constraints)
		a.cursor.Replace(n)

	case nodes.ConvertRowtypeExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.CopyStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Query", nil, n.Query)
		a.apply(&n, "Attlist", nil, n.Attlist)
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.CreateAmStmt:
		a.apply(&n, "HandlerName", nil, n.HandlerName)
		a.cursor.Replace(n)

	case nodes.CreateCastStmt:
		if n.Sourcetype != nil {
			a.apply(&n, "Sourcetype", nil, *n.Sourcetype)
		}
		if n.Targettype != nil {
			a.apply(&n, "Targettype", nil, *n.Targettype)
		}
		a.apply(&n, "Func", nil, n.Func)
		a.cursor.Replace(n)

	case nodes.CreateConversionStmt:
		a.apply(&n, "ConversionName", nil, n.ConversionName)
		a.apply(&n, "Funcname", nil, n.FuncName)
		a.cursor.Replace(n)

	case nodes.CreateDomainStmt:
		a.apply(&n, "Domainname", nil, n.Domainname)
		if n.TypeName != nil {
			a.apply(&n, "TypeName", nil, *n.TypeName)
		}
		if n.CollClause != nil {
			a.apply(&n, "CollClause", nil, *n.CollClause)
		}
		a.apply(&n, "Constraints", nil, n.Constraints)
		a.cursor.Replace(n)

	case nodes.CreateEnumStmt:
		a.apply(&n, "TypeName", nil, n.TypeName)
		a.apply(&n, "Vals", nil, n.Vals)
		a.cursor.Replace(n)

	case nodes.CreateEventTrigStmt:
		a.apply(&n, "", nil, n.Whenclause)
		a.apply(&n, "", nil, n.Funcname)

	case nodes.CreateExtensionStmt:
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateFdwStmt:
		a.apply(&n, "", nil, n.FuncOptions)
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateForeignServerStmt:
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateForeignTableStmt:
		a.apply(&n, "", nil, n.Base)
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateFunctionStmt:
		a.apply(&n, "", nil, n.Funcname)
		a.apply(&n, "", nil, n.Parameters)
		if n.ReturnType != nil {
			a.apply(&n, "", nil, *n.ReturnType)
		}
		a.apply(&n, "", nil, n.Options)
		a.apply(&n, "", nil, n.WithClause)

	case nodes.CreateOpClassItem:
		a.apply(&n, "", nil, n.Name)
		a.apply(&n, "", nil, n.OrderFamily)
		a.apply(&n, "", nil, n.ClassArgs)
		if n.Storedtype != nil {
			a.apply(&n, "", nil, *n.Storedtype)
		}

	case nodes.CreateOpClassStmt:
		a.apply(&n, "", nil, n.Opclassname)
		a.apply(&n, "", nil, n.Opfamilyname)
		if n.Datatype != nil {
			a.apply(&n, "", nil, *n.Datatype)
		}
		a.apply(&n, "", nil, n.Items)

	case nodes.CreateOpFamilyStmt:
		a.apply(&n, "", nil, n.Opfamilyname)

	case nodes.CreatePLangStmt:
		a.apply(&n, "", nil, n.Plhandler)
		a.apply(&n, "", nil, n.Plinline)
		a.apply(&n, "", nil, n.Plvalidator)

	case nodes.CreatePolicyStmt:
		if n.Table != nil {
			a.apply(&n, "", nil, *n.Table)
		}
		a.apply(&n, "", nil, n.Roles)
		a.apply(&n, "", nil, n.Qual)
		a.apply(&n, "", nil, n.WithCheck)

	case nodes.CreatePublicationStmt:
		a.apply(&n, "", nil, n.Options)
		a.apply(&n, "", nil, n.Tables)

	case nodes.CreateRangeStmt:
		a.apply(&n, "", nil, n.TypeName)
		a.apply(&n, "", nil, n.Params)

	case nodes.CreateRoleStmt:
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateSchemaStmt:
		if n.Authrole != nil {
			a.apply(&n, "", nil, *n.Authrole)
		}
		a.apply(&n, "", nil, n.SchemaElts)

	case nodes.CreateSeqStmt:
		if n.Sequence != nil {
			a.apply(&n, "", nil, *n.Sequence)
		}
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateStatsStmt:
		a.apply(&n, "", nil, n.Defnames)
		a.apply(&n, "", nil, n.StatTypes)
		a.apply(&n, "", nil, n.Exprs)
		a.apply(&n, "", nil, n.Relations)

	case nodes.CreateStmt:
		if n.Relation != nil {
			a.apply(&n, "", nil, *n.Relation)
		}
		a.apply(&n, "", nil, n.TableElts)
		a.apply(&n, "", nil, n.InhRelations)
		if n.Partbound != nil {
			a.apply(&n, "", nil, *n.Partbound)
		}
		if n.Partspec != nil {
			a.apply(&n, "", nil, *n.Partspec)
		}
		a.apply(&n, "", nil, n.Constraints)
		a.apply(&n, "", nil, n.Options)
		if n.OfTypename != nil {
			a.apply(&n, "", nil, *n.OfTypename)
		}

	case nodes.CreateSubscriptionStmt:
		a.apply(&n, "", nil, n.Publication)
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateTableAsStmt:
		a.apply(&n, "", nil, n.Query)
		a.apply(&n, "", nil, n.Into)

	case nodes.CreateTableSpaceStmt:
		if n.Owner != nil {
			a.apply(&n, "", nil, *n.Owner)
		}
		a.apply(&n, "", nil, n.Options)

	case nodes.CreateTransformStmt:
		if n.TypeName != nil {
			a.apply(&n, "", nil, *n.TypeName)
		}
		if n.Fromsql != nil {
			a.apply(&n, "", nil, *n.Fromsql)
		}
		if n.Tosql != nil {
			a.apply(&n, "", nil, *n.Tosql)
		}

	case nodes.CreateTrigStmt:
		if n.Relation != nil {
			a.apply(&n, "", nil, *n.Relation)
		}
		a.apply(&n, "", nil, n.Funcname)
		a.apply(&n, "", nil, n.Args)
		a.apply(&n, "", nil, n.Columns)
		a.apply(&n, "", nil, n.WhenClause)
		a.apply(&n, "", nil, n.TransitionRels)
		if n.Constrrel != nil {
			a.apply(&n, "", nil, *n.Constrrel)
		}

	case nodes.CreateUserMappingStmt:
		if n.User != nil {
			a.apply(&n, "", nil, *n.User)
		}
		a.apply(&n, "", nil, n.Options)

	case nodes.CreatedbStmt:
		a.apply(&n, "", nil, n.Options)

	case nodes.CurrentOfExpr:
		a.apply(&n, "", nil, n.Xpr)

	case nodes.DeallocateStmt:
		// pass

	case nodes.DeclareCursorStmt:
		a.apply(&n, "", nil, n.Query)

	case nodes.DefElem:
		a.apply(&n, "", nil, n.Arg)

	case nodes.DefineStmt:
		a.apply(&n, "", nil, n.Defnames)
		a.apply(&n, "", nil, n.Args)
		a.apply(&n, "", nil, n.Definition)

	case nodes.DeleteStmt:
		if n.Relation != nil {
			a.apply(&n, "", nil, *n.Relation)
		}
		a.apply(&n, "", nil, n.UsingClause)
		a.apply(&n, "", nil, n.WhereClause)
		a.apply(&n, "", nil, n.ReturningList)
		if n.WithClause != nil {
			a.apply(&n, "", nil, *n.WithClause)
		}

	case nodes.DiscardStmt:
		// pass

	case nodes.DoStmt:
		a.apply(&n, "", nil, n.Args)

	case nodes.DropOwnedStmt:
		a.apply(&n, "", nil, n.Roles)

	case nodes.DropRoleStmt:
		a.apply(&n, "", nil, n.Roles)

	case nodes.DropStmt:
		a.apply(&n, "", nil, n.Objects)

	case nodes.DropSubscriptionStmt:
		// pass

	case nodes.DropTableSpaceStmt:
		// pass

	case nodes.DropUserMappingStmt:
		if n.User != nil {
			a.apply(&n, "", nil, *n.User)
		}

	case nodes.DropdbStmt:
		// pass

	case nodes.ExecuteStmt:
		a.apply(&n, "", nil, n.Params)

	case nodes.ExplainStmt:
		a.apply(&n, "", nil, n.Query)
		a.apply(&n, "", nil, n.Options)

	case nodes.Expr:
		// pass

	case nodes.FetchStmt:
		// pass

	case nodes.FieldSelect:
		a.apply(&n, "", nil, n.Xpr)
		a.apply(&n, "", nil, n.Arg)

	case nodes.FieldStore:
		a.apply(&n, "", nil, n.Xpr)
		a.apply(&n, "", nil, n.Arg)
		a.apply(&n, "", nil, n.Newvals)
		a.apply(&n, "", nil, n.Fieldnums)

	case nodes.Float:
		// pass

	case nodes.FromExpr:
		a.apply(&n, "", nil, n.Fromlist)
		a.apply(&n, "", nil, n.Quals)

	case nodes.FuncCall:
		a.apply(&n, "", nil, n.Funcname)
		a.apply(&n, "", nil, n.Args)
		a.apply(&n, "", nil, n.AggOrder)
		a.apply(&n, "", nil, n.AggFilter)
		if n.Over != nil {
			a.apply(&n, "", nil, *n.Over)
		}

	case nodes.FuncExpr:
		a.apply(&n, "", nil, n.Xpr)
		a.apply(&n, "", nil, n.Args)

	case nodes.FunctionParameter:
		if n.ArgType != nil {
			a.apply(&n, "", nil, *n.ArgType)
		}
		a.apply(&n, "", nil, n.Defexpr)

	case nodes.GrantRoleStmt:
		a.apply(&n, "", nil, n.GrantedRoles)
		a.apply(&n, "", nil, n.GranteeRoles)
		if n.Grantor != nil {
			a.apply(&n, "", nil, *n.Grantor)
		}

	case nodes.GrantStmt:
		a.apply(&n, "", nil, n.Objects)
		a.apply(&n, "", nil, n.Privileges)
		a.apply(&n, "", nil, n.Grantees)

	case nodes.GroupingFunc:
		a.apply(&n, "", nil, n.Xpr)
		a.apply(&n, "", nil, n.Args)
		a.apply(&n, "", nil, n.Refs)
		a.apply(&n, "", nil, n.Cols)
		a.cursor.Replace(n)

	case nodes.GroupingSet:
		a.apply(&n, "Content", nil, n.Content)
		a.cursor.Replace(n)

	case nodes.ImportForeignSchemaStmt:
		a.apply(&n, "TableList", nil, n.TableList)
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.IndexElem:
		a.apply(&n, "Expr", nil, n.Expr)
		a.apply(&n, "Collation", nil, n.Collation)
		a.apply(&n, "Opclass", nil, n.Opclass)
		a.cursor.Replace(n)

	case nodes.IndexStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "IndexParams", nil, n.IndexParams)
		a.apply(&n, "Options", nil, n.Options)
		a.apply(&n, "WhereClause", nil, n.WhereClause)
		a.apply(&n, "ExcludeOpNames", nil, n.ExcludeOpNames)
		a.cursor.Replace(n)

	case nodes.InferClause:
		a.apply(&n, "IndexElems", nil, n.IndexElems)
		a.apply(&n, "WhereClause", nil, n.WhereClause)
		a.cursor.Replace(n)

	case nodes.InferenceElem:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Expr", nil, n.Expr)
		a.cursor.Replace(n)

	case nodes.InlineCodeBlock:
		// pass

	case nodes.InsertStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Cols", nil, n.Cols)
		a.apply(&n, "SelectStmt", nil, n.SelectStmt)
		if n.OnConflictClause != nil {
			a.apply(&n, "OnConflictClause", nil, *n.OnConflictClause)
		}
		a.apply(&n, "ReturningList", nil, n.ReturningList)
		if n.WithClause != nil {
			a.apply(&n, "WithClause", nil, *n.WithClause)
		}
		a.cursor.Replace(n)

	case nodes.Integer:
		// pass

	case nodes.IntoClause:
		if n.Rel != nil {
			a.apply(&n, "Rel", nil, *n.Rel)
		}
		a.apply(&n, "ColNames", nil, n.ColNames)
		a.apply(&n, "Options", nil, n.Options)
		a.apply(&n, "ViewQuery", nil, n.ViewQuery)
		a.cursor.Replace(n)

	case nodes.JoinExpr:
		a.apply(&n, "Larg", nil, n.Larg)
		a.apply(&n, "Rarg", nil, n.Rarg)
		a.apply(&n, "UsingClause", nil, n.UsingClause)
		a.apply(&n, "Quals", nil, n.Quals)
		if n.Alias != nil {
			a.apply(&n, "Alias", nil, *n.Alias)
		}
		a.cursor.Replace(n)

	case nodes.List:
		a.applyList(&n, "Items")
		a.cursor.Replace(n)

	case nodes.ListenStmt:
		// pass

	case nodes.LoadStmt:
		// pass

	case nodes.LockStmt:
		a.apply(&n, "Relations", nil, n.Relations)

	case nodes.LockingClause:
		a.apply(&n, "LockedRels", nil, n.LockedRels)

	case nodes.MinMaxExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Args)

	case nodes.MultiAssignRef:
		a.apply(&n, "Source", nil, n.Source)

	case nodes.NamedArgExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Arg)

	case nodes.NextValueExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)

	case nodes.NotifyStmt:
		// pass

	case nodes.Null:
		// pass

	case nodes.NullTest:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)

	case nodes.ObjectWithArgs:
		a.apply(&n, "Objname", nil, n.Objname)
		a.apply(&n, "Objargs", nil, n.Objargs)

	case nodes.OnConflictClause:
		if n.Infer != nil {
			a.apply(&n, "Infer", nil, *n.Infer)
		}
		a.apply(&n, "TargetList", nil, n.TargetList)
		a.apply(&n, "WhereClause", nil, n.WhereClause)

	case nodes.OnConflictExpr:
		a.apply(&n, "", nil, n.ArbiterElems)
		a.apply(&n, "", nil, n.ArbiterWhere)
		a.apply(&n, "", nil, n.OnConflictSet)
		a.apply(&n, "", nil, n.OnConflictWhere)
		a.apply(&n, "", nil, n.ExclRelTlist)

	case nodes.OpExpr:
		a.apply(&n, "", nil, n.Xpr)
		a.apply(&n, "", nil, n.Args)

	case nodes.Param:
		a.apply(&n, "", nil, n.Xpr)

	case nodes.ParamExecData:
		// pass

	case nodes.ParamExternData:
		// pass

	case nodes.ParamListInfoData:
		// pass

	case nodes.ParamRef:
		// pass

	case nodes.PartitionBoundSpec:
		a.apply(&n, "", nil, n.Listdatums)
		a.apply(&n, "", nil, n.Lowerdatums)
		a.apply(&n, "", nil, n.Upperdatums)

	case nodes.PartitionCmd:
		if n.Name != nil {
			a.apply(&n, "", nil, *n.Name)
		}
		if n.Bound != nil {
			a.apply(&n, "", nil, *n.Bound)
		}

	case nodes.PartitionElem:
		a.apply(&n, "", nil, n.Expr)
		a.apply(&n, "", nil, n.Collation)
		a.apply(&n, "", nil, n.Opclass)

	case nodes.PartitionRangeDatum:
		a.apply(&n, "", nil, n.Value)

	case nodes.PartitionSpec:
		a.apply(&n, "", nil, n.PartParams)

	case nodes.PrepareStmt:
		a.apply(&n, "", nil, n.Argtypes)
		a.apply(&n, "", nil, n.Query)

	case nodes.Query:
		a.apply(&n, "", nil, n.UtilityStmt)
		a.apply(&n, "", nil, n.CteList)
		a.apply(&n, "", nil, n.Jointree)
		a.apply(&n, "", nil, n.TargetList)
		a.apply(&n, "", nil, n.OnConflict)
		a.apply(&n, "", nil, n.ReturningList)
		a.apply(&n, "", nil, n.GroupClause)
		a.apply(&n, "", nil, n.GroupingSets)
		a.apply(&n, "", nil, n.HavingQual)
		a.apply(&n, "", nil, n.WindowClause)
		a.apply(&n, "", nil, n.DistinctClause)
		a.apply(&n, "", nil, n.SortClause)
		a.apply(&n, "", nil, n.LimitCount)
		a.apply(&n, "", nil, n.RowMarks)
		a.apply(&n, "", nil, n.SetOperations)
		a.apply(&n, "", nil, n.ConstraintDeps)
		a.apply(&n, "", nil, n.WithCheckOptions)

	case nodes.RangeFunction:
		a.apply(&n, "", nil, n.Functions)
		if n.Alias != nil {
			a.apply(&n, "", nil, *n.Alias)
		}
		a.apply(&n, "", nil, n.Coldeflist)

	case nodes.RangeSubselect:
		a.apply(&n, "", nil, n.Subquery)
		if n.Alias != nil {
			a.apply(&n, "", nil, *n.Alias)
		}

	case nodes.RangeTableFunc:
		a.apply(&n, "", nil, n.Docexpr)
		a.apply(&n, "", nil, n.Rowexpr)
		a.apply(&n, "", nil, n.Namespaces)
		a.apply(&n, "", nil, n.Columns)
		if n.Alias != nil {
			a.apply(&n, "", nil, *n.Alias)
		}

	case nodes.RangeTableFuncCol:
		if n.TypeName != nil {
			a.apply(&n, "", nil, *n.TypeName)
		}
		a.apply(&n, "", nil, n.Colexpr)
		a.apply(&n, "", nil, n.Coldefexpr)

	case nodes.RangeTableSample:
		a.apply(&n, "", nil, n.Relation)
		a.apply(&n, "", nil, n.Method)
		a.apply(&n, "", nil, n.Args)

	case nodes.RangeTblEntry:
		a.apply(&n, "", nil, n.Tablesample)
		a.apply(&n, "", nil, n.Subquery)
		a.apply(&n, "", nil, n.Joinaliasvars)
		a.apply(&n, "", nil, n.Functions)
		a.apply(&n, "", nil, n.Tablefunc)
		a.apply(&n, "", nil, n.ValuesLists)
		a.apply(&n, "", nil, n.Coltypes)
		a.apply(&n, "", nil, n.Colcollations)
		if n.Alias != nil {
			a.apply(&n, "", nil, *n.Alias)
		}
		a.apply(&n, "", nil, n.Eref)
		a.apply(&n, "", nil, n.SecurityQuals)

	case nodes.RangeTblFunction:
		a.apply(&n, "", nil, n.Funcexpr)
		a.apply(&n, "", nil, n.Funccolnames)
		a.apply(&n, "", nil, n.Funccoltypes)
		a.apply(&n, "", nil, n.Funccoltypmods)
		a.apply(&n, "", nil, n.Funccolcollations)

	case nodes.RangeTblRef:
		// pass

	case nodes.RangeVar:
		if n.Alias != nil {
			a.apply(&n, "Alias", nil, *n.Alias)
			a.cursor.Replace(n)
		}

	case nodes.RawStmt:
		a.apply(&n, "Stmt", nil, n.Stmt)
		a.cursor.Replace(n)

	case nodes.ReassignOwnedStmt:
		a.apply(&n, "Roles", nil, n.Roles)
		if n.Newrole != nil {
			a.apply(&n, "Newrole", nil, *n.Newrole)
		}
		a.cursor.Replace(n)

	case nodes.RefreshMatViewStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
			a.cursor.Replace(n)
		}

	case nodes.ReindexStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
			a.cursor.Replace(n)
		}

	case nodes.RelabelType:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Arg", nil, n.Arg)
		a.cursor.Replace(n)

	case nodes.RenameStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "Object", nil, n.Object)
		a.cursor.Replace(n)

	case nodes.ReplicaIdentityStmt:
		// pass

	case nodes.ResTarget:
		a.apply(&n, "Indirection", nil, n.Indirection)
		a.apply(&n, "Val", nil, n.Val)
		a.cursor.Replace(n)

	case nodes.RoleSpec:
		// pass

	case nodes.RowCompareExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Opnos", nil, n.Opnos)
		a.apply(&n, "Opfamilies", nil, n.Opfamilies)
		a.apply(&n, "Inputcollids", nil, n.Inputcollids)
		a.apply(&n, "Largs", nil, n.Largs)
		a.apply(&n, "Rargs", nil, n.Rargs)
		a.cursor.Replace(n)

	case nodes.RowExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Args)
		a.apply(&n, "Colnames", nil, n.Colnames)
		a.cursor.Replace(n)

	case nodes.RowMarkClause:
		// pass

	case nodes.RuleStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "WhereClause", nil, n.WhereClause)
		a.apply(&n, "Actions", nil, n.Actions)
		a.cursor.Replace(n)

	case nodes.SQLValueFunction:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.cursor.Replace(n)

	case nodes.ScalarArrayOpExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Args)
		a.cursor.Replace(n)

	case nodes.SecLabelStmt:
		a.apply(&n, "Object", nil, n.Object)
		a.cursor.Replace(n)

	case nodes.SelectStmt:
		a.apply(&n, "DistinctClause", nil, n.DistinctClause)
		if n.IntoClause != nil {
			a.apply(&n, "IntoClause", nil, *n.IntoClause)
		}
		a.apply(&n, "TargetList", nil, n.TargetList)
		a.apply(&n, "FromClause", nil, n.FromClause)
		a.apply(&n, "WhereClause", nil, n.WhereClause)
		a.apply(&n, "GroupClause", nil, n.GroupClause)
		a.apply(&n, "HavingClause", nil, n.HavingClause)
		a.apply(&n, "WindowClause", nil, n.WindowClause)
		// TODO: Not sure how to handle a slice of a slice
		//
		// for _, vs := range n.ValuesLists {
		// 	for _, v := range vs {
		// 		a.apply(&n, "", nil, v)
		// 	}
		// }
		a.apply(&n, "SortClause", nil, n.SortClause)
		a.apply(&n, "LimitOffset", nil, n.LimitOffset)
		a.apply(&n, "LimitCount", nil, n.LimitCount)
		a.apply(&n, "LockingClause", nil, n.LockingClause)
		if n.WithClause != nil {
			a.apply(&n, "WithClause", nil, *n.WithClause)
		}
		if n.Larg != nil {
			a.apply(&n, "Larg", nil, *n.Larg)
		}
		if n.Rarg != nil {
			a.apply(&n, "Rarg", nil, *n.Rarg)
		}
		a.cursor.Replace(n)

	case nodes.SetOperationStmt:
		a.apply(&n, "Larg", nil, n.Larg)
		a.apply(&n, "Rarg", nil, n.Rarg)
		a.apply(&n, "ColTypes", nil, n.ColTypes)
		a.apply(&n, "ColTypmods", nil, n.ColTypmods)
		a.apply(&n, "ColCollations", nil, n.ColCollations)
		a.apply(&n, "GroupClauses", nil, n.GroupClauses)
		a.cursor.Replace(n)

	case nodes.SetToDefault:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.cursor.Replace(n)

	case nodes.SortBy:
		a.apply(&n, "Node", nil, n.Node)
		a.apply(&n, "UseOp", nil, n.UseOp)
		a.cursor.Replace(n)

	case nodes.SortGroupClause:
		// pass

	case nodes.String:
		// pass

	case nodes.SubLink:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Testexpr", nil, n.Testexpr)
		a.apply(&n, "Opername", nil, n.OperName)
		a.apply(&n, "Subselect", nil, n.Subselect)
		a.cursor.Replace(n)

	case nodes.SubPlan:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Testexpr", nil, n.Testexpr)
		a.apply(&n, "ParamIds", nil, n.ParamIds)
		a.apply(&n, "SetParam", nil, n.SetParam)
		a.apply(&n, "ParParam", nil, n.ParParam)
		a.apply(&n, "Args", nil, n.Args)
		a.cursor.Replace(n)

	case nodes.TableFunc:
		a.apply(&n, "NsUris", nil, n.NsUris)
		a.apply(&n, "NsNames", nil, n.NsNames)
		a.apply(&n, "Docexpr", nil, n.Docexpr)
		a.apply(&n, "Rowexpr", nil, n.Rowexpr)
		a.apply(&n, "Colnames", nil, n.Colnames)
		a.apply(&n, "Coltypes", nil, n.Coltypes)
		a.apply(&n, "ColTypmods", nil, n.Coltypmods)
		a.apply(&n, "Colcollations", nil, n.Colcollations)
		a.apply(&n, "Colexprs", nil, n.Colexprs)
		a.apply(&n, "Coldefexprs", nil, n.Coldefexprs)
		a.cursor.Replace(n)

	case nodes.TableLikeClause:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
			a.cursor.Replace(n)
		}

	case nodes.TableSampleClause:
		a.apply(&n, "Args", nil, n.Args)
		a.apply(&n, "Repeatable", nil, n.Repeatable)
		a.cursor.Replace(n)

	case nodes.TargetEntry:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Expr", nil, n.Expr)
		a.cursor.Replace(n)

	case nodes.TransactionStmt:
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.TriggerTransition:
		// pass

	case nodes.TruncateStmt:
		a.apply(&n, "Relations", nil, n.Relations)
		a.cursor.Replace(n)

	case nodes.TypeCast:
		a.apply(&n, "Arg", nil, n.Arg)
		a.apply(&n, "TypeName", nil, n.TypeName)
		a.cursor.Replace(n)

	case nodes.TypeName:
		a.apply(&n, "Names", nil, n.Names)
		a.apply(&n, "Typmods", nil, n.Typmods)
		a.apply(&n, "ArrayBounds", nil, n.ArrayBounds)
		a.cursor.Replace(n)

	case *nodes.TypeName:
		a.apply(n, "Names", nil, n.Names)
		a.apply(n, "Typmods", nil, n.Typmods)
		a.apply(n, "ArrayBounds", nil, n.ArrayBounds)
		a.cursor.Replace(n)

	case nodes.UnlistenStmt:
		// pass

	case nodes.UpdateStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "TargetList", nil, n.TargetList)
		a.apply(&n, "WhereClause", nil, n.WhereClause)
		a.apply(&n, "FromClause", nil, n.FromClause)
		a.apply(&n, "ReturningList", nil, n.ReturningList)
		if n.WithClause != nil {
			a.apply(&n, "WithClause", nil, *n.WithClause)
		}
		a.cursor.Replace(n)

	case nodes.VacuumStmt:
		if n.Relation != nil {
			a.apply(&n, "Relation", nil, *n.Relation)
		}
		a.apply(&n, "VaCols", nil, n.VaCols)
		a.cursor.Replace(n)

	case nodes.Var:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.cursor.Replace(n)

	case nodes.VariableSetStmt:
		a.apply(&n, "Args", nil, n.Args)
		a.cursor.Replace(n)

	case nodes.VariableShowStmt:
		// pass

	case nodes.ViewStmt:
		if n.View != nil {
			a.apply(&n, "View", nil, *n.View)
		}
		a.apply(&n, "Aliases", nil, n.Aliases)
		a.apply(&n, "Query", nil, n.Query)
		a.apply(&n, "Options", nil, n.Options)
		a.cursor.Replace(n)

	case nodes.WindowClause:
		a.apply(&n, "PartitionClause", nil, n.PartitionClause)
		a.apply(&n, "OrderClause", nil, n.OrderClause)
		a.apply(&n, "StartOffset", nil, n.StartOffset)
		a.apply(&n, "EndOffset", nil, n.EndOffset)
		a.cursor.Replace(n)

	case nodes.WindowDef:
		a.apply(&n, "PartitionClause", nil, n.PartitionClause)
		a.apply(&n, "OrderClause", nil, n.OrderClause)
		a.apply(&n, "StartOffset", nil, n.StartOffset)
		a.apply(&n, "EndOffset", nil, n.EndOffset)
		a.cursor.Replace(n)

	case nodes.WindowFunc:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "Args", nil, n.Args)
		a.apply(&n, "Aggfilter", nil, n.Aggfilter)
		a.cursor.Replace(n)

	case nodes.WithCheckOption:
		a.apply(&n, "Qual", nil, n.Qual)
		a.cursor.Replace(n)

	case nodes.WithClause:
		a.apply(&n, "Ctes", nil, n.Ctes)
		a.cursor.Replace(n)

	case nodes.XmlExpr:
		a.apply(&n, "Xpr", nil, n.Xpr)
		a.apply(&n, "NamedArgs", nil, n.NamedArgs)
		a.apply(&n, "ArgNames", nil, n.ArgNames)
		a.apply(&n, "Args", nil, n.Args)
		a.cursor.Replace(n)

	case nodes.XmlSerialize:
		a.apply(&n, "Expr", nil, n.Expr)
		if n.TypeName != nil {
			a.apply(&n, "TypeName", nil, *n.TypeName)
		}
		a.cursor.Replace(n)

	default:
		panic(fmt.Sprintf("Apply: unexpected node type %T", n))
	}

	if a.post != nil && !a.post(&a.cursor) {
		panic(abort)
	}

	a.cursor = saved
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
