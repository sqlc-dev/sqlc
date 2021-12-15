package printer

import (
	"strconv"
	"strings"

	"github.com/kyleconroy/sqlc/internal/python/ast"
)

type writer struct {
	options Options
	src     []byte
}

type Options struct {
}

type PrintResult struct {
	Python []byte
}

func Print(node *ast.Node, options Options) PrintResult {
	w := writer{options: options}
	w.printNode(node, 0)
	return PrintResult{
		Python: w.src,
	}
}

func (w *writer) print(text string) {
	w.src = append(w.src, text...)
}

func (w *writer) printIndent(indent int32) {
	for i, n := 0, int(indent); i < n; i++ {
		w.src = append(w.src, "    "...)
	}
}

func (w *writer) printNode(node *ast.Node, indent int32) {
	switch n := node.Node.(type) {

	case *ast.Node_Alias:
		w.print(n.Alias.Name)

	case *ast.Node_AnnAssign:
		w.printAnnAssign(n.AnnAssign, indent)

	case *ast.Node_Assign:
		w.printAssign(n.Assign, indent)

	case *ast.Node_AsyncFor:
		w.printAsyncFor(n.AsyncFor, indent)

	case *ast.Node_AsyncFunctionDef:
		w.printAsyncFunctionDef(n.AsyncFunctionDef, indent)

	case *ast.Node_Attribute:
		w.printAttribute(n.Attribute, indent)

	case *ast.Node_Await:
		w.printAwait(n.Await, indent)

	case *ast.Node_Call:
		w.printCall(n.Call, indent)

	case *ast.Node_ClassDef:
		w.printClassDef(n.ClassDef, indent)

	case *ast.Node_Comment:
		w.printComment(n.Comment, indent)

	case *ast.Node_Compare:
		w.printCompare(n.Compare, indent)

	case *ast.Node_Constant:
		w.printConstant(n.Constant, indent)

	case *ast.Node_Dict:
		w.printDict(n.Dict, indent)

	case *ast.Node_Expr:
		w.printNode(n.Expr.Value, indent)

	case *ast.Node_For:
		w.printFor(n.For, indent)

	case *ast.Node_FunctionDef:
		w.printFunctionDef(n.FunctionDef, indent)

	case *ast.Node_If:
		w.printIf(n.If, indent)

	case *ast.Node_Import:
		w.printImport(n.Import, indent)

	case *ast.Node_ImportFrom:
		w.printImportFrom(n.ImportFrom, indent)

	case *ast.Node_ImportGroup:
		w.printImportGroup(n.ImportGroup, indent)

	case *ast.Node_Is:
		w.print("is")

	case *ast.Node_Keyword:
		w.printKeyword(n.Keyword, indent)

	case *ast.Node_Module:
		w.printModule(n.Module, indent)

	case *ast.Node_Name:
		w.print(n.Name.Id)

	case *ast.Node_Pass:
		w.print("pass")

	case *ast.Node_Return:
		w.printReturn(n.Return, indent)

	case *ast.Node_Subscript:
		w.printSubscript(n.Subscript, indent)

	case *ast.Node_Yield:
		w.printYield(n.Yield, indent)

	default:
		panic(n)

	}
}

func (w *writer) printAnnAssign(aa *ast.AnnAssign, indent int32) {
	if aa.Comment != "" {
		w.print("# ")
		w.print(aa.Comment)
		w.print("\n")
		w.printIndent(indent)
	}
	w.printName(aa.Target, indent)
	w.print(": ")
	w.printNode(aa.Annotation, indent)
}

func (w *writer) printArg(a *ast.Arg, indent int32) {
	w.print(a.Arg)
	if a.Annotation != nil {
		w.print(": ")
		w.printNode(a.Annotation, indent)
	}
}

func (w *writer) printAssign(a *ast.Assign, indent int32) {
	for i, name := range a.Targets {
		w.printNode(name, indent)
		if i != len(a.Targets)-1 {
			w.print(", ")
		}
	}
	w.print(" = ")
	w.printNode(a.Value, indent)
}

func (w *writer) printAsyncFor(n *ast.AsyncFor, indent int32) {
	w.print("async ")
	w.printFor(&ast.For{
		Target: n.Target,
		Iter:   n.Iter,
		Body:   n.Body,
	}, indent)
}

func (w *writer) printAsyncFunctionDef(afd *ast.AsyncFunctionDef, indent int32) {
	w.print("async ")
	w.printFunctionDef(&ast.FunctionDef{
		Name:    afd.Name,
		Args:    afd.Args,
		Body:    afd.Body,
		Returns: afd.Returns,
	}, indent)
}

func (w *writer) printAttribute(a *ast.Attribute, indent int32) {
	if _, ok := a.Value.Node.(*ast.Node_Await); ok {
		w.print("(")
		w.printNode(a.Value, indent)
		w.print(")")
	} else {
		w.printNode(a.Value, indent)
	}
	w.print(".")
	w.print(a.Attr)
}

func (w *writer) printAwait(n *ast.Await, indent int32) {
	w.print("await ")
	w.printNode(n.Value, indent)
}

func (w *writer) printCall(c *ast.Call, indent int32) {
	w.printNode(c.Func, indent)
	w.print("(")
	for i, a := range c.Args {
		w.printNode(a, indent)
		if i != len(c.Args)-1 {
			w.print(", ")
		}
	}
	for _, kw := range c.Keywords {
		w.print("\n")
		w.printIndent(indent + 1)
		w.printKeyword(kw, indent+1)
		w.print(",")
	}
	if len(c.Keywords) > 0 {
		w.print("\n")
		w.printIndent(indent)
	}
	w.print(")")
}

func (w *writer) printClassDef(cd *ast.ClassDef, indent int32) {
	for _, node := range cd.DecoratorList {
		w.print("@")
		w.printNode(node, indent)
		w.print("\n")
	}
	w.print("class ")
	w.print(cd.Name)
	if len(cd.Bases) > 0 {
		w.print("(")
		for i, node := range cd.Bases {
			w.printNode(node, indent)
			if i != len(cd.Bases)-1 {
				w.print(", ")
			}
		}
		w.print(")")
	}
	w.print(":\n")
	for i, node := range cd.Body {
		if i != 0 {
			if _, ok := node.Node.(*ast.Node_FunctionDef); ok {
				w.print("\n")
			}
			if _, ok := node.Node.(*ast.Node_AsyncFunctionDef); ok {
				w.print("\n")
			}
		}
		w.printIndent(indent + 1)
		// A docstring is a string literal that occurs as the first
		// statement in a module, function, class, or method
		// definition. Such a docstring becomes the __doc__ special
		// attribute of that object.
		if i == 0 {
			if e, ok := node.Node.(*ast.Node_Expr); ok {
				if c, ok := e.Expr.Value.Node.(*ast.Node_Constant); ok {
					w.print(`""`)
					w.printConstant(c.Constant, indent)
					w.print(`""`)
					w.print("\n")
					continue
				}
			}
		}
		w.printNode(node, indent+1)
		w.print("\n")
	}
}

func (w *writer) printConstant(c *ast.Constant, indent int32) {
	switch n := c.Value.(type) {
	case *ast.Constant_Int:
		w.print(strconv.Itoa(int(n.Int)))

	case *ast.Constant_None:
		w.print("None")

	case *ast.Constant_Str:
		str := `"`
		if strings.Contains(n.Str, "\n") {
			str = `"""`
		}
		w.print(str)
		w.print(n.Str)
		w.print(str)

	default:
		panic(n)
	}
}

func (w *writer) printComment(c *ast.Comment, indent int32) {
	w.print("# ")
	w.print(c.Text)
	w.print("\n")
}

func (w *writer) printCompare(c *ast.Compare, indent int32) {
	w.printNode(c.Left, indent)
	w.print(" ")
	for _, node := range c.Ops {
		w.printNode(node, indent)
		w.print(" ")
	}
	for _, node := range c.Comparators {
		w.printNode(node, indent)
	}
}

func (w *writer) printDict(d *ast.Dict, indent int32) {
	if len(d.Keys) != len(d.Values) {
		panic(`dict keys and values are not the same length`)
	}
	w.print("{")
	split := len(d.Keys) > 3
	keyIndent := indent
	if split {
		keyIndent += 1
	}
	for i, _ := range d.Keys {
		if split {
			w.print("\n")
			w.printIndent(keyIndent)
		}
		w.printNode(d.Keys[i], keyIndent)
		w.print(": ")
		w.printNode(d.Values[i], keyIndent)
		if i != len(d.Keys)-1 || split {
			if split {
				w.print(",")
			} else {
				w.print(", ")
			}
		}
	}
	if split {
		w.print("\n")
		w.printIndent(indent)
	}
	w.print("}")
}

func (w *writer) printFor(n *ast.For, indent int32) {
	w.print("for ")
	w.printNode(n.Target, indent)
	w.print(" in ")
	w.printNode(n.Iter, indent)
	w.print(":\n")
	for i, node := range n.Body {
		w.printIndent(indent + 1)
		w.printNode(node, indent+1)
		if i != len(n.Body)-1 {
			w.print("\n")
		}
	}
}

func (w *writer) printIf(i *ast.If, indent int32) {
	w.print("if ")
	w.printNode(i.Test, indent)
	w.print(":\n")
	for j, node := range i.Body {
		w.printIndent(indent + 1)
		w.printNode(node, indent+1)
		if j != len(i.Body)-1 {
			w.print("\n")
		}
	}
}

func (w *writer) printFunctionDef(fd *ast.FunctionDef, indent int32) {
	w.print("def ")
	w.print(fd.Name)
	w.print("(")
	if fd.Args != nil {
		for i, arg := range fd.Args.Args {
			w.printArg(arg, indent)
			if i != len(fd.Args.Args)-1 {
				w.print(", ")
			}
		}
		if len(fd.Args.KwOnlyArgs) > 0 {
			w.print(", *, ")
			for i, arg := range fd.Args.KwOnlyArgs {
				w.printArg(arg, indent)
				if i != len(fd.Args.KwOnlyArgs)-1 {
					w.print(", ")
				}
			}
		}
	}
	w.print(")")
	if fd.Returns != nil {
		w.print(" -> ")
		w.printNode(fd.Returns, indent)
	}
	w.print(":\n")
	for i, node := range fd.Body {
		w.printIndent(indent + 1)
		w.printNode(node, indent+1)
		if i != len(fd.Body)-1 {
			w.print("\n")
		}
	}
}

func (w *writer) printImport(imp *ast.Import, indent int32) {
	w.print("import ")
	for i, node := range imp.Names {
		w.printNode(node, indent)
		if i != len(imp.Names)-1 {
			w.print(", ")
		}
	}
	w.print("\n")
}

func (w *writer) printImportFrom(imp *ast.ImportFrom, indent int32) {
	w.print("from ")
	w.print(imp.Module)
	w.print(" import ")
	for i, node := range imp.Names {
		w.printNode(node, indent)
		if i != len(imp.Names)-1 {
			w.print(", ")
		}
	}
	w.print("\n")
}

func (w *writer) printImportGroup(n *ast.ImportGroup, indent int32) {
	if len(n.Imports) == 0 {
		return
	}
	for _, node := range n.Imports {
		w.printNode(node, indent)
	}
	w.print("\n")
}

func (w *writer) printIs(i *ast.Is, indent int32) {
	w.print("is")
}
func (w *writer) printKeyword(k *ast.Keyword, indent int32) {
	w.print(k.Arg)
	w.print("=")
	w.printNode(k.Value, indent)
}

func (w *writer) printModule(mod *ast.Module, indent int32) {
	for i, node := range mod.Body {
		prevIsImport := false
		if i > 0 {
			_, isImport := mod.Body[i-1].Node.(*ast.Node_ImportGroup)
			prevIsImport = isImport
		}
		_, isClassDef := node.Node.(*ast.Node_ClassDef)
		_, isAssign := node.Node.(*ast.Node_Assign)
		if isClassDef || isAssign {
			if prevIsImport {
				w.print("\n")
			} else {
				w.print("\n\n")
			}
		}
		w.printNode(node, indent)
		if isAssign {
			w.print("\n")
		}
	}
}

func (w *writer) printName(n *ast.Name, indent int32) {
	w.print(n.Id)
}

func (w *writer) printReturn(r *ast.Return, indent int32) {
	w.print("return ")
	w.printNode(r.Value, indent)
}

func (w *writer) printSubscript(ss *ast.Subscript, indent int32) {
	w.printName(ss.Value, indent)
	w.print("[")
	w.printNode(ss.Slice, indent)
	w.print("]")

}

func (w *writer) printYield(n *ast.Yield, indent int32) {
	w.print("yield ")
	w.printNode(n.Value, indent)
}
