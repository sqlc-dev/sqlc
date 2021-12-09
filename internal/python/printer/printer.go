package printer

import "github.com/kyleconroy/sqlc/internal/python/ast"

type writer struct {
	options Options
	src     []byte
}

type Options struct {
}

type PrintResult struct {
	Code []byte
}

func Print(node *ast.Node, options Options) PrintResult {
	w := writer{options: options}
	w.printNode(node, 0)
	return PrintResult{
		Code: w.src,
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

	case *ast.Node_ClassDef:
		w.printClassDef(n.ClassDef, indent)

	case *ast.Node_Import:
		w.printImport(n.Import, indent)

	case *ast.Node_ImportFrom:
		w.printImportFrom(n.ImportFrom, indent)

	case *ast.Node_Module:
		w.printModule(n.Module, indent)

	case *ast.Node_Name:
		w.print(n.Name.Id)

	default:
		panic(n)

	}
}

func (w *writer) printAnnAssign(aa *ast.AnnAssign, indent int32) {
	w.print(aa.Target.Id)
	w.print(": ")
	w.printNode(aa.Annotation, indent)
}

func (w *writer) printClassDef(cd *ast.ClassDef, indent int32) {
	for _, node := range cd.DecoratorList {
		w.print("@")
		w.printNode(node, indent)
		w.print("\n")
	}
	w.print("class ")
	w.print(cd.Name)
	w.print(":\n")
	for _, node := range cd.Body {
		w.printIndent(indent + 1)
		w.printNode(node, indent+1)
		w.print("\n")
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
		w.printNode(node, indent)
	}
}
