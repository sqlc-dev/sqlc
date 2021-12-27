package printer

import (
	"github.com/kyleconroy/sqlc/internal/kotlin/ast"
)

type writer struct {
	options Options
	src     []byte
}

type Options struct {
}

type PrintResult struct {
	Kotlin []byte
}

func Print(node *ast.Node, options Options) PrintResult {
	w := writer{options: options}
	w.printNode(node, 0)
	return PrintResult{
		Kotlin: w.src,
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

	case *ast.Node_Class:
		w.printClass(n.Class, indent)

	case *ast.Node_Comment:
		w.printComment(n.Comment, indent)

	case *ast.Node_DeclarationModifierList:
		w.printDeclarationModifierList(n.DeclarationModifierList, indent)

	case *ast.Node_DotQualifiedExpression:
		w.printDotQualifiedExpression(n.DotQualifiedExpression, indent)

	case *ast.Node_File:
		w.printFile(n.File, indent)

	case *ast.Node_NameReferenceExpression:
		w.printNameReferenceExpression(n.NameReferenceExpression, indent)

	case *ast.Node_PackageDirective:
		w.printPackageDirective(n.PackageDirective, indent)

	default:
		panic(n)

	}
}

func (w *writer) printClass(n *ast.Class, indent int32) {
	if n.Modifiers != nil {
		w.printDeclarationModifierList(n.Modifiers, indent)
	}
	w.print("class ")
	w.print(n.Name)
	w.print(" (\n")
	for _, node := range n.Body {
		w.printNode(node, indent+1)
	}
	w.print(")\n")
}

func (w *writer) printComment(n *ast.Comment, indent int32) {
	w.print("// ")
	w.print(n.Text)
	w.print("\n")
}

func (w *writer) printDeclarationModifierList(n *ast.DeclarationModifierList, indent int32) {
	for _, a := range n.Annotations {
		w.print(a)
		w.print(" ")
	}
}

func (w *writer) printDotQualifiedExpression(n *ast.DotQualifiedExpression, indent int32) {
	w.printNode(n.Receiver, indent)
	w.print(".")
	w.printNode(n.Selector, indent)
}

func (w *writer) printFile(n *ast.File, indent int32) {
	for _, node := range n.Body {
		_, isClass := node.Node.(*ast.Node_Class)
		if isClass {
			w.print("\n")
		}
		w.printNode(node, indent)
	}
}

func (w *writer) printNameReferenceExpression(n *ast.NameReferenceExpression, indent int32) {
	w.print(n.Name)
}

func (w *writer) printPackageDirective(n *ast.PackageDirective, indent int32) {
	w.print("\n")
	w.print("package ")
	w.printNode(n.Name, indent)
	w.print("\n")
}
