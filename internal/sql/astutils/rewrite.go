// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package astutils

import (
	"fmt"
	"reflect"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
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
func Apply(root ast.Node, pre, post ApplyFunc) (result ast.Node) {
	parent := &struct{ ast.Node }{root}
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
	parent ast.Node
	name   string
	iter   *iterator // valid if non-nil
	node   ast.Node
}

// Node returns the current Node.
func (c *Cursor) Node() ast.Node { return c.node }

// Parent returns the parent of the current Node.
func (c *Cursor) Parent() ast.Node { return c.parent }

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
func (c *Cursor) Replace(n ast.Node) {
	v := c.field()
	if i := c.Index(); i >= 0 {
		v = v.Index(i)
	}
	v.Set(reflect.ValueOf(n))
}

// D// application carries all the shared data so we can pass it around cheaply.
type application struct {
	pre, post ApplyFunc
	cursor    Cursor
	iter      iterator
}

func (a *application) apply(parent ast.Node, name string, iter *iterator, n ast.Node) {
	// convert typed nil into untyped nil
	if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
		n = nil
	}

	// avoid heap-allocating a new cursor for each apply call; reuse a.cursor instead
	saved := a.cursor
	a.cursor.parent = parent
	a.cursor.name = name
	a.cursor.iter = iter
	a.cursor.node = n

	if a.pre != nil && !a.pre(&a.cursor) {
		a.cursor = saved
		return
	}

	// walk children
	// (the order of the cases matches the order of the corresponding node types in go/ast)
	switch n := n.(type) {
	case nil:
		// nothing to do

	case *ast.AlterTableCmd:
		a.apply(n, "Def", nil, n.Def)

	case *ast.AlterTableSetSchemaStmt:
		a.apply(n, "Table", nil, n.Table)

	case *ast.AlterTableStmt:
		a.apply(n, "Table", nil, n.Table)
		a.apply(n, "Cmds", nil, n.Cmds)

	case *ast.AlterTypeAddValueStmt:
		a.apply(n, "Type", nil, n.Type)

	case *ast.AlterTypeRenameValueStmt:
		a.apply(n, "Type", nil, n.Type)

	case *ast.ColumnDef:
		a.apply(n, "TypeName", nil, n.TypeName)

	case *ast.ColumnRef:
		// pass

	case *ast.CommentOnColumnStmt:
		a.apply(n, "Table", nil, n.Table)
		a.apply(n, "Col", nil, n.Col)

	case *ast.CommentOnSchemaStmt:
		a.apply(n, "Schema", nil, n.Schema)

	case *ast.CommentOnTableStmt:
		a.apply(n, "Table", nil, n.Table)

	case *ast.CommentOnTypeStmt:
		a.apply(n, "Type", nil, n.Type)

	case *ast.CreateEnumStmt:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "Vals", nil, n.Vals)

	case *ast.CreateFunctionStmt:
		a.apply(n, "ReturnType", nil, n.ReturnType)
		a.apply(n, "Func", nil, n.Func)

	case *ast.CreateSchemaStmt:
		// pass

	case *ast.CreateTableStmt:
		a.apply(n, "Name", nil, n.Name)

	case *ast.DropFunctionStmt:
		// pass

	case *ast.DropSchemaStmt:
		// pass

	case *ast.DropTableStmt:
		// pass

	case *ast.DropTypeStmt:
		// pass

	case *ast.FuncName:
		// pass

	case *ast.FuncParam:
		a.apply(n, "Type", nil, n.Type)
		a.apply(n, "DefExpr", nil, n.DefExpr)

	case *ast.FuncSpec:
		a.apply(n, "Name", nil, n.Name)

	case *ast.List:
		// Since item is a slice
		a.applyList(n, "Items")

	case *ast.RawStmt:
		a.apply(n, "Stmt", nil, n.Stmt)

	case *ast.RenameColumnStmt:
		a.apply(n, "Table", nil, n.Table)
		a.apply(n, "Col", nil, n.Col)

	case *ast.RenameTableStmt:
		a.apply(n, "Table", nil, n.Table)

	case *ast.ResTarget:
		a.apply(n, "Val", nil, n.Val)

	case *ast.SelectStmt:
		a.apply(n, "Fields", nil, n.Fields)
		a.apply(n, "From", nil, n.From)

	case *ast.Statement:
		a.apply(n, "Raw", nil, n.Raw)

	case *ast.String:
		// pass

	case *ast.TODO:
		// pass

	case *ast.TableName:
		// pass

	case *ast.TypeName:
		// pass

	case *pg.A_ArrayExpr:
		a.apply(n, "Elements", nil, n.Elements)

	case *pg.A_Const:
		a.apply(n, "Val", nil, n.Val)

	case *pg.A_Expr:
		a.apply(n, "Name", nil, n.Name)
		a.apply(n, "Lexpr", nil, n.Lexpr)
		a.apply(n, "Rexpr", nil, n.Rexpr)

	case *pg.A_Indices:
		a.apply(n, "Lidx", nil, n.Lidx)
		a.apply(n, "Uidx", nil, n.Uidx)

	case *pg.A_Indirection:
		a.apply(n, "Arg", nil, n.Arg)
		a.apply(n, "Indirection", nil, n.Indirection)

	case *pg.A_Star:
		// pass

	case *pg.AccessPriv:
		a.apply(n, "Cols", nil, n.Cols)

	case *pg.Aggref:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Aggargtypes", nil, n.Aggargtypes)
		a.apply(n, "Aggdirectargs", nil, n.Aggdirectargs)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Aggorder", nil, n.Aggorder)
		a.apply(n, "Aggdistinct", nil, n.Aggdistinct)
		a.apply(n, "Aggfilter", nil, n.Aggfilter)

	case *pg.Alias:
		a.apply(n, "Colnames", nil, n.Colnames)

	case *pg.AlterCollationStmt:
		a.apply(n, "Collname", nil, n.Collname)

	case *pg.AlterDatabaseSetStmt:
		a.apply(n, "Setstmt", nil, n.Setstmt)

	case *pg.AlterDatabaseStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterDefaultPrivilegesStmt:
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "Action", nil, n.Action)

	case *pg.AlterDomainStmt:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "Def", nil, n.Def)

	case *pg.AlterEnumStmt:
		a.apply(n, "TypeName", nil, n.TypeName)

	case *pg.AlterEventTrigStmt:
		// pass

	case *pg.AlterExtensionContentsStmt:
		a.apply(n, "Object", nil, n.Object)

	case *pg.AlterExtensionStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterFdwStmt:
		a.apply(n, "FuncOptions", nil, n.FuncOptions)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterForeignServerStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterFunctionStmt:
		a.apply(n, "Func", nil, n.Func)
		a.apply(n, "Actions", nil, n.Actions)

	case *pg.AlterObjectDependsStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Object", nil, n.Object)
		a.apply(n, "Extname", nil, n.Extname)

	case *pg.AlterObjectSchemaStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Object", nil, n.Object)

	case *pg.AlterOpFamilyStmt:
		a.apply(n, "Opfamilyname", nil, n.Opfamilyname)
		a.apply(n, "Items", nil, n.Items)

	case *pg.AlterOperatorStmt:
		a.apply(n, "Opername", nil, n.Opername)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterOwnerStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Object", nil, n.Object)
		a.apply(n, "Newowner", nil, n.Newowner)

	case *pg.AlterPolicyStmt:
		a.apply(n, "Table", nil, n.Table)
		a.apply(n, "Roles", nil, n.Roles)
		a.apply(n, "Qual", nil, n.Qual)
		a.apply(n, "WithCheck", nil, n.WithCheck)

	case *pg.AlterPublicationStmt:
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "Tables", nil, n.Tables)

	case *pg.AlterRoleSetStmt:
		a.apply(n, "Role", nil, n.Role)
		a.apply(n, "Setstmt", nil, n.Setstmt)

	case *pg.AlterRoleStmt:
		a.apply(n, "Role", nil, n.Role)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterSeqStmt:
		a.apply(n, "Sequence", nil, n.Sequence)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterSubscriptionStmt:
		a.apply(n, "Publication", nil, n.Publication)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterSystemStmt:
		a.apply(n, "Setstmt", nil, n.Setstmt)

	case *pg.AlterTSConfigurationStmt:
		a.apply(n, "Cfgname", nil, n.Cfgname)
		a.apply(n, "Tokentype", nil, n.Tokentype)
		a.apply(n, "Dicts", nil, n.Dicts)

	case *pg.AlterTSDictionaryStmt:
		a.apply(n, "Dictname", nil, n.Dictname)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterTableCmd:
		a.apply(n, "Newowner", nil, n.Newowner)
		a.apply(n, "Def", nil, n.Def)

	case *pg.AlterTableMoveAllStmt:
		a.apply(n, "Roles", nil, n.Roles)

	case *pg.AlterTableSpaceOptionsStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlterTableStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Cmds", nil, n.Cmds)

	case *pg.AlterUserMappingStmt:
		a.apply(n, "User", nil, n.User)
		a.apply(n, "Options", nil, n.Options)

	case *pg.AlternativeSubPlan:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Subplans", nil, n.Subplans)

	case *pg.ArrayCoerceExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.ArrayExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Elements", nil, n.Elements)

	case *pg.ArrayRef:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Refupperindexpr", nil, n.Refupperindexpr)
		a.apply(n, "Reflowerindexpr", nil, n.Reflowerindexpr)
		a.apply(n, "Refexpr", nil, n.Refexpr)
		a.apply(n, "Refassgnexpr", nil, n.Refassgnexpr)

	case *pg.BitString:
		// pass

	case *pg.BlockIdData:
		// pass

	case *pg.BoolExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)

	case *pg.BooleanTest:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.CaseExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Defresult", nil, n.Defresult)

	case *pg.CaseTestExpr:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.CaseWhen:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Expr", nil, n.Expr)
		a.apply(n, "Result", nil, n.Result)

	case *pg.CheckPointStmt:
		// pass

	case *pg.ClosePortalStmt:
		// pass

	case *pg.ClusterStmt:
		a.apply(n, "Relation", nil, n.Relation)

	case *pg.CoalesceExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)

	case *pg.CoerceToDomain:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.CoerceToDomainValue:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.CoerceViaIO:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.CollateClause:
		a.apply(n, "Arg", nil, n.Arg)
		a.apply(n, "Collname", nil, n.Collname)

	case *pg.CollateExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.ColumnDef:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "RawDefault", nil, n.RawDefault)
		a.apply(n, "CookedDefault", nil, n.CookedDefault)
		a.apply(n, "CollClause", nil, n.CollClause)
		a.apply(n, "Constraints", nil, n.Constraints)
		a.apply(n, "Fdwoptions", nil, n.Fdwoptions)

	case *pg.ColumnRef:
		a.apply(n, "Fields", nil, n.Fields)

	case *pg.CommentStmt:
		a.apply(n, "Object", nil, n.Object)

	case *pg.CommonTableExpr:
		a.apply(n, "Aliascolnames", nil, n.Aliascolnames)
		a.apply(n, "Ctequery", nil, n.Ctequery)
		a.apply(n, "Ctecolnames", nil, n.Ctecolnames)
		a.apply(n, "Ctecoltypes", nil, n.Ctecoltypes)
		a.apply(n, "Ctecoltypmods", nil, n.Ctecoltypmods)
		a.apply(n, "Ctecolcollations", nil, n.Ctecolcollations)

	case *pg.CompositeTypeStmt:
		a.apply(n, "Typevar", nil, n.Typevar)
		a.apply(n, "Coldeflist", nil, n.Coldeflist)

	case *pg.Const:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.Constraint:
		a.apply(n, "RawExpr", nil, n.RawExpr)
		a.apply(n, "Keys", nil, n.Keys)
		a.apply(n, "Exclusions", nil, n.Exclusions)
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "WhereClause", nil, n.WhereClause)
		a.apply(n, "Pktable", nil, n.Pktable)
		a.apply(n, "FkAttrs", nil, n.FkAttrs)
		a.apply(n, "PkAttrs", nil, n.PkAttrs)
		a.apply(n, "OldConpfeqop", nil, n.OldConpfeqop)

	case *pg.ConstraintsSetStmt:
		a.apply(n, "Constraints", nil, n.Constraints)

	case *pg.ConvertRowtypeExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.CopyStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Query", nil, n.Query)
		a.apply(n, "Attlist", nil, n.Attlist)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateAmStmt:
		a.apply(n, "HandlerName", nil, n.HandlerName)

	case *pg.CreateCastStmt:
		a.apply(n, "Sourcetype", nil, n.Sourcetype)
		a.apply(n, "Targettype", nil, n.Targettype)
		a.apply(n, "Func", nil, n.Func)

	case *pg.CreateConversionStmt:
		a.apply(n, "ConversionName", nil, n.ConversionName)
		a.apply(n, "FuncName", nil, n.FuncName)

	case *pg.CreateDomainStmt:
		a.apply(n, "Domainname", nil, n.Domainname)
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "CollClause", nil, n.CollClause)
		a.apply(n, "Constraints", nil, n.Constraints)

	case *pg.CreateEnumStmt:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "Vals", nil, n.Vals)

	case *pg.CreateEventTrigStmt:
		a.apply(n, "Whenclause", nil, n.Whenclause)
		a.apply(n, "Funcname", nil, n.Funcname)

	case *pg.CreateExtensionStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateFdwStmt:
		a.apply(n, "FuncOptions", nil, n.FuncOptions)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateForeignServerStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateForeignTableStmt:
		a.apply(n, "Base", nil, n.Base)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateFunctionStmt:
		a.apply(n, "Funcname", nil, n.Funcname)
		a.apply(n, "Parameters", nil, n.Parameters)
		a.apply(n, "ReturnType", nil, n.ReturnType)
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "WithClause", nil, n.WithClause)

	case *pg.CreateOpClassItem:
		a.apply(n, "Name", nil, n.Name)
		a.apply(n, "OrderFamily", nil, n.OrderFamily)
		a.apply(n, "ClassArgs", nil, n.ClassArgs)
		a.apply(n, "Storedtype", nil, n.Storedtype)

	case *pg.CreateOpClassStmt:
		a.apply(n, "Opclassname", nil, n.Opclassname)
		a.apply(n, "Opfamilyname", nil, n.Opfamilyname)
		a.apply(n, "Datatype", nil, n.Datatype)
		a.apply(n, "Items", nil, n.Items)

	case *pg.CreateOpFamilyStmt:
		a.apply(n, "Opfamilyname", nil, n.Opfamilyname)

	case *pg.CreatePLangStmt:
		a.apply(n, "Plhandler", nil, n.Plhandler)
		a.apply(n, "Plinline", nil, n.Plinline)
		a.apply(n, "Plvalidator", nil, n.Plvalidator)

	case *pg.CreatePolicyStmt:
		a.apply(n, "Table", nil, n.Table)
		a.apply(n, "Roles", nil, n.Roles)
		a.apply(n, "Qual", nil, n.Qual)
		a.apply(n, "WithCheck", nil, n.WithCheck)

	case *pg.CreatePublicationStmt:
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "Tables", nil, n.Tables)

	case *pg.CreateRangeStmt:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "Params", nil, n.Params)

	case *pg.CreateRoleStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateSchemaStmt:
		a.apply(n, "Authrole", nil, n.Authrole)
		a.apply(n, "SchemaElts", nil, n.SchemaElts)

	case *pg.CreateSeqStmt:
		a.apply(n, "Sequence", nil, n.Sequence)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateStatsStmt:
		a.apply(n, "Defnames", nil, n.Defnames)
		a.apply(n, "StatTypes", nil, n.StatTypes)
		a.apply(n, "Exprs", nil, n.Exprs)
		a.apply(n, "Relations", nil, n.Relations)

	case *pg.CreateStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "TableElts", nil, n.TableElts)
		a.apply(n, "InhRelations", nil, n.InhRelations)
		a.apply(n, "Partbound", nil, n.Partbound)
		a.apply(n, "Partspec", nil, n.Partspec)
		a.apply(n, "OfTypename", nil, n.OfTypename)
		a.apply(n, "Constraints", nil, n.Constraints)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateSubscriptionStmt:
		a.apply(n, "Publication", nil, n.Publication)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateTableAsStmt:
		a.apply(n, "Query", nil, n.Query)
		a.apply(n, "Into", nil, n.Into)

	case *pg.CreateTableSpaceStmt:
		a.apply(n, "Owner", nil, n.Owner)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreateTransformStmt:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "Fromsql", nil, n.Fromsql)
		a.apply(n, "Tosql", nil, n.Tosql)

	case *pg.CreateTrigStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Funcname", nil, n.Funcname)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Columns", nil, n.Columns)
		a.apply(n, "WhenClause", nil, n.WhenClause)
		a.apply(n, "TransitionRels", nil, n.TransitionRels)
		a.apply(n, "Constrrel", nil, n.Constrrel)

	case *pg.CreateUserMappingStmt:
		a.apply(n, "User", nil, n.User)
		a.apply(n, "Options", nil, n.Options)

	case *pg.CreatedbStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.CurrentOfExpr:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.DeallocateStmt:
		// pass

	case *pg.DeclareCursorStmt:
		a.apply(n, "Query", nil, n.Query)

	case *pg.DefElem:
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.DefineStmt:
		a.apply(n, "Defnames", nil, n.Defnames)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Definition", nil, n.Definition)

	case *pg.DeleteStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "UsingClause", nil, n.UsingClause)
		a.apply(n, "WhereClause", nil, n.WhereClause)
		a.apply(n, "ReturningList", nil, n.ReturningList)
		a.apply(n, "WithClause", nil, n.WithClause)

	case *pg.DiscardStmt:
		// pass

	case *pg.DoStmt:
		a.apply(n, "Args", nil, n.Args)

	case *pg.DropOwnedStmt:
		a.apply(n, "Roles", nil, n.Roles)

	case *pg.DropRoleStmt:
		a.apply(n, "Roles", nil, n.Roles)

	case *pg.DropStmt:
		a.apply(n, "Objects", nil, n.Objects)

	case *pg.DropSubscriptionStmt:
		// pass

	case *pg.DropTableSpaceStmt:
		// pass

	case *pg.DropUserMappingStmt:
		a.apply(n, "User", nil, n.User)

	case *pg.DropdbStmt:
		// pass

	case *pg.ExecuteStmt:
		a.apply(n, "Params", nil, n.Params)

	case *pg.ExplainStmt:
		a.apply(n, "Query", nil, n.Query)
		a.apply(n, "Options", nil, n.Options)

	case *pg.Expr:
		// pass

	case *pg.FetchStmt:
		// pass

	case *pg.FieldSelect:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.FieldStore:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)
		a.apply(n, "Newvals", nil, n.Newvals)
		a.apply(n, "Fieldnums", nil, n.Fieldnums)

	case *pg.Float:
		// pass

	case *pg.FromExpr:
		a.apply(n, "Fromlist", nil, n.Fromlist)
		a.apply(n, "Quals", nil, n.Quals)

	case *pg.FuncCall:
		a.apply(n, "Funcname", nil, n.Funcname)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "AggOrder", nil, n.AggOrder)
		a.apply(n, "AggFilter", nil, n.AggFilter)
		a.apply(n, "Over", nil, n.Over)

	case *pg.FuncExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)

	case *pg.FunctionParameter:
		a.apply(n, "ArgType", nil, n.ArgType)
		a.apply(n, "Defexpr", nil, n.Defexpr)

	case *pg.GrantRoleStmt:
		a.apply(n, "GrantedRoles", nil, n.GrantedRoles)
		a.apply(n, "GranteeRoles", nil, n.GranteeRoles)
		a.apply(n, "Grantor", nil, n.Grantor)

	case *pg.GrantStmt:
		a.apply(n, "Objects", nil, n.Objects)
		a.apply(n, "Privileges", nil, n.Privileges)
		a.apply(n, "Grantees", nil, n.Grantees)

	case *pg.GroupingFunc:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Refs", nil, n.Refs)
		a.apply(n, "Cols", nil, n.Cols)

	case *pg.GroupingSet:
		a.apply(n, "Content", nil, n.Content)

	case *pg.ImportForeignSchemaStmt:
		a.apply(n, "TableList", nil, n.TableList)
		a.apply(n, "Options", nil, n.Options)

	case *pg.IndexElem:
		a.apply(n, "Expr", nil, n.Expr)
		a.apply(n, "Collation", nil, n.Collation)
		a.apply(n, "Opclass", nil, n.Opclass)

	case *pg.IndexStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "IndexParams", nil, n.IndexParams)
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "WhereClause", nil, n.WhereClause)
		a.apply(n, "ExcludeOpNames", nil, n.ExcludeOpNames)

	case *pg.InferClause:
		a.apply(n, "IndexElems", nil, n.IndexElems)
		a.apply(n, "WhereClause", nil, n.WhereClause)

	case *pg.InferenceElem:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Expr", nil, n.Expr)

	case *pg.InlineCodeBlock:
		// pass

	case *pg.InsertStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Cols", nil, n.Cols)
		a.apply(n, "SelectStmt", nil, n.SelectStmt)
		a.apply(n, "OnConflictClause", nil, n.OnConflictClause)
		a.apply(n, "ReturningList", nil, n.ReturningList)
		a.apply(n, "WithClause", nil, n.WithClause)

	case *pg.Integer:
		// pass

	case *pg.IntoClause:
		a.apply(n, "Rel", nil, n.Rel)
		a.apply(n, "ColNames", nil, n.ColNames)
		a.apply(n, "Options", nil, n.Options)
		a.apply(n, "ViewQuery", nil, n.ViewQuery)

	case *pg.JoinExpr:
		a.apply(n, "Larg", nil, n.Larg)
		a.apply(n, "Rarg", nil, n.Rarg)
		a.apply(n, "UsingClause", nil, n.UsingClause)
		a.apply(n, "Quals", nil, n.Quals)
		a.apply(n, "Alias", nil, n.Alias)

	case *pg.ListenStmt:
		// pass

	case *pg.LoadStmt:
		// pass

	case *pg.LockStmt:
		a.apply(n, "Relations", nil, n.Relations)

	case *pg.LockingClause:
		a.apply(n, "LockedRels", nil, n.LockedRels)

	case *pg.MinMaxExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)

	case *pg.MultiAssignRef:
		a.apply(n, "Source", nil, n.Source)

	case *pg.NamedArgExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.NextValueExpr:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.NotifyStmt:
		// pass

	case *pg.Null:
		// pass

	case *pg.NullTest:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.ObjectWithArgs:
		a.apply(n, "Objname", nil, n.Objname)
		a.apply(n, "Objargs", nil, n.Objargs)

	case *pg.OnConflictClause:
		a.apply(n, "Infer", nil, n.Infer)
		a.apply(n, "TargetList", nil, n.TargetList)
		a.apply(n, "WhereClause", nil, n.WhereClause)

	case *pg.OnConflictExpr:
		a.apply(n, "ArbiterElems", nil, n.ArbiterElems)
		a.apply(n, "ArbiterWhere", nil, n.ArbiterWhere)
		a.apply(n, "OnConflictSet", nil, n.OnConflictSet)
		a.apply(n, "OnConflictWhere", nil, n.OnConflictWhere)
		a.apply(n, "ExclRelTlist", nil, n.ExclRelTlist)

	case *pg.OpExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)

	case *pg.Param:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.ParamExecData:
		// pass

	case *pg.ParamExternData:
		// pass

	case *pg.ParamListInfoData:
		// pass

	case *pg.ParamRef:
		// pass

	case *pg.PartitionBoundSpec:
		a.apply(n, "Listdatums", nil, n.Listdatums)
		a.apply(n, "Lowerdatums", nil, n.Lowerdatums)
		a.apply(n, "Upperdatums", nil, n.Upperdatums)

	case *pg.PartitionCmd:
		a.apply(n, "Name", nil, n.Name)
		a.apply(n, "Bound", nil, n.Bound)

	case *pg.PartitionElem:
		a.apply(n, "Expr", nil, n.Expr)
		a.apply(n, "Collation", nil, n.Collation)
		a.apply(n, "Opclass", nil, n.Opclass)

	case *pg.PartitionRangeDatum:
		a.apply(n, "Value", nil, n.Value)

	case *pg.PartitionSpec:
		a.apply(n, "PartParams", nil, n.PartParams)

	case *pg.PrepareStmt:
		a.apply(n, "Argtypes", nil, n.Argtypes)
		a.apply(n, "Query", nil, n.Query)

	case *pg.Query:
		a.apply(n, "UtilityStmt", nil, n.UtilityStmt)
		a.apply(n, "CteList", nil, n.CteList)
		a.apply(n, "Rtable", nil, n.Rtable)
		a.apply(n, "Jointree", nil, n.Jointree)
		a.apply(n, "TargetList", nil, n.TargetList)
		a.apply(n, "OnConflict", nil, n.OnConflict)
		a.apply(n, "ReturningList", nil, n.ReturningList)
		a.apply(n, "GroupClause", nil, n.GroupClause)
		a.apply(n, "GroupingSets", nil, n.GroupingSets)
		a.apply(n, "HavingQual", nil, n.HavingQual)
		a.apply(n, "WindowClause", nil, n.WindowClause)
		a.apply(n, "DistinctClause", nil, n.DistinctClause)
		a.apply(n, "SortClause", nil, n.SortClause)
		a.apply(n, "LimitOffset", nil, n.LimitOffset)
		a.apply(n, "LimitCount", nil, n.LimitCount)
		a.apply(n, "RowMarks", nil, n.RowMarks)
		a.apply(n, "SetOperations", nil, n.SetOperations)
		a.apply(n, "ConstraintDeps", nil, n.ConstraintDeps)
		a.apply(n, "WithCheckOptions", nil, n.WithCheckOptions)

	case *pg.RangeFunction:
		a.apply(n, "Functions", nil, n.Functions)
		a.apply(n, "Alias", nil, n.Alias)
		a.apply(n, "Coldeflist", nil, n.Coldeflist)

	case *pg.RangeSubselect:
		a.apply(n, "Subquery", nil, n.Subquery)
		a.apply(n, "Alias", nil, n.Alias)

	case *pg.RangeTableFunc:
		a.apply(n, "Docexpr", nil, n.Docexpr)
		a.apply(n, "Rowexpr", nil, n.Rowexpr)
		a.apply(n, "Namespaces", nil, n.Namespaces)
		a.apply(n, "Columns", nil, n.Columns)
		a.apply(n, "Alias", nil, n.Alias)

	case *pg.RangeTableFuncCol:
		a.apply(n, "TypeName", nil, n.TypeName)
		a.apply(n, "Colexpr", nil, n.Colexpr)
		a.apply(n, "Coldefexpr", nil, n.Coldefexpr)

	case *pg.RangeTableSample:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Method", nil, n.Method)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Repeatable", nil, n.Repeatable)

	case *pg.RangeTblEntry:
		a.apply(n, "Tablesample", nil, n.Tablesample)
		a.apply(n, "Subquery", nil, n.Subquery)
		a.apply(n, "Joinaliasvars", nil, n.Joinaliasvars)
		a.apply(n, "Functions", nil, n.Functions)
		a.apply(n, "Tablefunc", nil, n.Tablefunc)
		a.apply(n, "ValuesLists", nil, n.ValuesLists)
		a.apply(n, "Coltypes", nil, n.Coltypes)
		a.apply(n, "Coltypmods", nil, n.Coltypmods)
		a.apply(n, "Colcollations", nil, n.Colcollations)
		a.apply(n, "Alias", nil, n.Alias)
		a.apply(n, "Eref", nil, n.Eref)
		a.apply(n, "SecurityQuals", nil, n.SecurityQuals)

	case *pg.RangeTblFunction:
		a.apply(n, "Funcexpr", nil, n.Funcexpr)
		a.apply(n, "Funccolnames", nil, n.Funccolnames)
		a.apply(n, "Funccoltypes", nil, n.Funccoltypes)
		a.apply(n, "Funccoltypmods", nil, n.Funccoltypmods)
		a.apply(n, "Funccolcollations", nil, n.Funccolcollations)

	case *pg.RangeTblRef:
		// pass

	case *pg.RangeVar:
		a.apply(n, "Alias", nil, n.Alias)

	case *pg.RawStmt:
		a.apply(n, "Stmt", nil, n.Stmt)

	case *pg.ReassignOwnedStmt:
		a.apply(n, "Roles", nil, n.Roles)
		a.apply(n, "Newrole", nil, n.Newrole)

	case *pg.RefreshMatViewStmt:
		a.apply(n, "Relation", nil, n.Relation)

	case *pg.ReindexStmt:
		a.apply(n, "Relation", nil, n.Relation)

	case *pg.RelabelType:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Arg", nil, n.Arg)

	case *pg.RenameStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "Object", nil, n.Object)

	case *pg.ReplicaIdentityStmt:
		// pass

	case *pg.ResTarget:
		a.apply(n, "Indirection", nil, n.Indirection)
		a.apply(n, "Val", nil, n.Val)

	case *pg.RoleSpec:
		// pass

	case *pg.RowCompareExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Opnos", nil, n.Opnos)
		a.apply(n, "Opfamilies", nil, n.Opfamilies)
		a.apply(n, "Inputcollids", nil, n.Inputcollids)
		a.apply(n, "Largs", nil, n.Largs)
		a.apply(n, "Rargs", nil, n.Rargs)

	case *pg.RowExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Colnames", nil, n.Colnames)

	case *pg.RowMarkClause:
		// pass

	case *pg.RuleStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "WhereClause", nil, n.WhereClause)
		a.apply(n, "Actions", nil, n.Actions)

	case *pg.SQLValueFunction:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.ScalarArrayOpExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)

	case *pg.SecLabelStmt:
		a.apply(n, "Object", nil, n.Object)

	case *pg.SelectStmt:
		a.apply(n, "DistinctClause", nil, n.DistinctClause)
		a.apply(n, "IntoClause", nil, n.IntoClause)
		a.apply(n, "TargetList", nil, n.TargetList)
		a.apply(n, "FromClause", nil, n.FromClause)
		a.apply(n, "WhereClause", nil, n.WhereClause)
		a.apply(n, "GroupClause", nil, n.GroupClause)
		a.apply(n, "HavingClause", nil, n.HavingClause)
		a.apply(n, "WindowClause", nil, n.WindowClause)
		a.apply(n, "SortClause", nil, n.SortClause)
		a.apply(n, "LimitOffset", nil, n.LimitOffset)
		a.apply(n, "LimitCount", nil, n.LimitCount)
		a.apply(n, "LockingClause", nil, n.LockingClause)
		a.apply(n, "WithClause", nil, n.WithClause)
		a.apply(n, "Larg", nil, n.Larg)
		a.apply(n, "Rarg", nil, n.Rarg)

	case *pg.SetOperationStmt:
		a.apply(n, "Larg", nil, n.Larg)
		a.apply(n, "Rarg", nil, n.Rarg)
		a.apply(n, "ColTypes", nil, n.ColTypes)
		a.apply(n, "ColTypmods", nil, n.ColTypmods)
		a.apply(n, "ColCollations", nil, n.ColCollations)
		a.apply(n, "GroupClauses", nil, n.GroupClauses)

	case *pg.SetToDefault:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.SortBy:
		a.apply(n, "Node", nil, n.Node)
		a.apply(n, "UseOp", nil, n.UseOp)

	case *pg.SortGroupClause:
		// pass

	case *pg.String:
		// pass

	case *pg.SubLink:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Testexpr", nil, n.Testexpr)
		a.apply(n, "OperName", nil, n.OperName)
		a.apply(n, "Subselect", nil, n.Subselect)

	case *pg.SubPlan:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Testexpr", nil, n.Testexpr)
		a.apply(n, "ParamIds", nil, n.ParamIds)
		a.apply(n, "SetParam", nil, n.SetParam)
		a.apply(n, "ParParam", nil, n.ParParam)
		a.apply(n, "Args", nil, n.Args)

	case *pg.TableFunc:
		a.apply(n, "NsUris", nil, n.NsUris)
		a.apply(n, "NsNames", nil, n.NsNames)
		a.apply(n, "Docexpr", nil, n.Docexpr)
		a.apply(n, "Rowexpr", nil, n.Rowexpr)
		a.apply(n, "Colnames", nil, n.Colnames)
		a.apply(n, "Coltypes", nil, n.Coltypes)
		a.apply(n, "Coltypmods", nil, n.Coltypmods)
		a.apply(n, "Colcollations", nil, n.Colcollations)
		a.apply(n, "Colexprs", nil, n.Colexprs)
		a.apply(n, "Coldefexprs", nil, n.Coldefexprs)

	case *pg.TableLikeClause:
		a.apply(n, "Relation", nil, n.Relation)

	case *pg.TableSampleClause:
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Repeatable", nil, n.Repeatable)

	case *pg.TargetEntry:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Expr", nil, n.Expr)

	case *pg.TransactionStmt:
		a.apply(n, "Options", nil, n.Options)

	case *pg.TriggerTransition:
		// pass

	case *pg.TruncateStmt:
		a.apply(n, "Relations", nil, n.Relations)

	case *pg.TypeCast:
		a.apply(n, "Arg", nil, n.Arg)
		a.apply(n, "TypeName", nil, n.TypeName)

	case *pg.TypeName:
		a.apply(n, "Names", nil, n.Names)
		a.apply(n, "Typmods", nil, n.Typmods)
		a.apply(n, "ArrayBounds", nil, n.ArrayBounds)

	case *pg.UnlistenStmt:
		// pass

	case *pg.UpdateStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "TargetList", nil, n.TargetList)
		a.apply(n, "WhereClause", nil, n.WhereClause)
		a.apply(n, "FromClause", nil, n.FromClause)
		a.apply(n, "ReturningList", nil, n.ReturningList)
		a.apply(n, "WithClause", nil, n.WithClause)

	case *pg.VacuumStmt:
		a.apply(n, "Relation", nil, n.Relation)
		a.apply(n, "VaCols", nil, n.VaCols)

	case *pg.Var:
		a.apply(n, "Xpr", nil, n.Xpr)

	case *pg.VariableSetStmt:
		a.apply(n, "Args", nil, n.Args)

	case *pg.VariableShowStmt:
		// pass

	case *pg.ViewStmt:
		a.apply(n, "View", nil, n.View)
		a.apply(n, "Aliases", nil, n.Aliases)
		a.apply(n, "Query", nil, n.Query)
		a.apply(n, "Options", nil, n.Options)

	case *pg.WindowClause:
		a.apply(n, "PartitionClause", nil, n.PartitionClause)
		a.apply(n, "OrderClause", nil, n.OrderClause)
		a.apply(n, "StartOffset", nil, n.StartOffset)
		a.apply(n, "EndOffset", nil, n.EndOffset)

	case *pg.WindowDef:
		a.apply(n, "PartitionClause", nil, n.PartitionClause)
		a.apply(n, "OrderClause", nil, n.OrderClause)
		a.apply(n, "StartOffset", nil, n.StartOffset)
		a.apply(n, "EndOffset", nil, n.EndOffset)

	case *pg.WindowFunc:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "Args", nil, n.Args)
		a.apply(n, "Aggfilter", nil, n.Aggfilter)

	case *pg.WithCheckOption:
		a.apply(n, "Qual", nil, n.Qual)

	case *pg.WithClause:
		a.apply(n, "Ctes", nil, n.Ctes)

	case *pg.XmlExpr:
		a.apply(n, "Xpr", nil, n.Xpr)
		a.apply(n, "NamedArgs", nil, n.NamedArgs)
		a.apply(n, "ArgNames", nil, n.ArgNames)
		a.apply(n, "Args", nil, n.Args)

	case *pg.XmlSerialize:
		a.apply(n, "Expr", nil, n.Expr)
		a.apply(n, "TypeName", nil, n.TypeName)

	// Comments and fields
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

func (a *application) applyList(parent ast.Node, name string) {
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
		var x ast.Node
		if e := v.Index(a.iter.index); e.IsValid() {
			x = e.Interface().(ast.Node)
		}

		a.iter.step = 1
		a.apply(parent, name, &a.iter, x)
		a.iter.index += a.iter.step
	}
	a.iter = saved
}
