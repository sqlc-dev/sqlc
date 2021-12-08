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

func (w *writer) printNode(node *ast.Node, indent int32) {
	switch n := node.Node.(type) {

	case *ast.Node_Alias:
		w.print(n.Alias.Name)

	case *ast.Node_ClassDef:
		w.printClassDef(n.ClassDef, indent)

	case *ast.Node_Import:
		w.printImport(n.Import, indent)

	case *ast.Node_ImportFrom:
		// w.printImport(n.Import, indent)

	case *ast.Node_Module:
		w.printModule(n.Module, indent)

	}
}

func (w *writer) printClassDef(cd *ast.ClassDef, indent int32) {
	w.print("\n\n")
	w.print("class ")
	w.print(cd.Name)
	w.print(":\n")
	w.print("    pass\n")
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

func (w *writer) printModule(mod *ast.Module, indent int32) {
	for _, node := range mod.Body {
		w.printNode(node, indent)
	}
}
