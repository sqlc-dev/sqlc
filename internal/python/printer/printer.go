package printer

import "github.com/kyleconroy/sqlc/internal/python/ast"

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

	case *ast.Node_Attribute:
		w.printAttribute(n.Attribute, indent)

	case *ast.Node_Call:
		w.printCall(n.Call, indent)

	case *ast.Node_ClassDef:
		w.printClassDef(n.ClassDef, indent)

	case *ast.Node_Comment:
		w.printComment(n.Comment, indent)

	case *ast.Node_Constant:
		w.printConstant(n.Constant, indent)

	case *ast.Node_Expr:
		w.printNode(n.Expr.Value, indent)

	case *ast.Node_Import:
		w.printImport(n.Import, indent)

	case *ast.Node_ImportFrom:
		w.printImportFrom(n.ImportFrom, indent)

	case *ast.Node_Module:
		w.printModule(n.Module, indent)

	case *ast.Node_Name:
		w.print(n.Name.Id)

	case *ast.Node_Subscript:
		w.printSubscript(n.Subscript, indent)

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

func (w *writer) printAssign(a *ast.Assign, indent int32) {
	for i, name := range a.Targets {
		w.printName(name, indent)
		if i != len(a.Targets)-1 {
			w.print(", ")
		}
	}
	w.print(" = ")
	w.printNode(a.Value, indent)
}

func (w *writer) printAttribute(a *ast.Attribute, indent int32) {
	w.printName(a.Value, indent)
	w.print(".")
	w.print(a.Attr)
}

func (w *writer) printCall(c *ast.Call, indent int32) {
	w.printNode(c.Func, indent)
	w.print("()")
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
		w.printIndent(indent + 1)
		// A docstring is a string literal that occurs as the first
		// statement in a module, function, class, or method
		// definition. Such a docstring becomes the __doc__ special
		// attribute of that object.
		if i == 0 {
			if e, ok := node.Node.(*ast.Node_Expr); ok {
				if c, ok := e.Expr.Value.Node.(*ast.Node_Constant); ok {
					w.print(`"""`)
					w.print(c.Constant.Value)
					w.print(`"""`)
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
	w.print("\"")
	w.print(c.Value)
	w.print("\"")
}

func (w *writer) printComment(c *ast.Comment, indent int32) {
	w.print("# ")
	w.print(c.Text)
	w.print("\n")
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

func (w *writer) printModule(mod *ast.Module, indent int32) {
	for _, node := range mod.Body {
		if _, ok := node.Node.(*ast.Node_ClassDef); ok {
			w.print("\n\n")
		}
		w.printNode(node, indent)
	}
}

func (w *writer) printName(n *ast.Name, indent int32) {
	w.print(n.Id)
}

func (w *writer) printSubscript(ss *ast.Subscript, indent int32) {
	w.printName(ss.Value, indent)
	w.print("[")
	w.printNode(ss.Slice, indent)
	w.print("]")

}
