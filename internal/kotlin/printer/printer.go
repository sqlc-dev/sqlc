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
		w.src = append(w.src, "  "...)
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
	if n.Constructor != nil {
		w.print(" ")
		w.printPrimaryConstructor(n.Constructor, indent)
	}
	for _, node := range n.Body {
		w.printNode(node, indent+1)
	}
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

func (w *writer) printParameter(n *ast.Parameter, indent int32) {
	w.print("val ")
	w.print(n.Name)
	w.print(": ")
	w.printTypeReference(n.Ref, indent)
}

func (w *writer) printParameterList(n *ast.ParameterList, indent int32) {
	for i, p := range n.Parameters {
		w.printIndent(indent)
		w.printParameter(p, indent)
		if i != len(n.Parameters)-1 {
			w.print(",")
		}
		w.print("\n")
	}
}

func (w *writer) printPrimaryConstructor(n *ast.PrimaryConstructor, indent int32) {
	w.print("(\n")
	if n.ValueParameterList != nil {
		w.printParameterList(n.ValueParameterList, indent+1)
	}
	w.print(")\n")
}

func (w *writer) printTypeReference(n *ast.TypeReference, indent int32) {
	w.printUserType(n.Element, indent)
}

func (w *writer) printUserType(n *ast.UserType, indent int32) {
	w.printNameReferenceExpression(n.ReferenceExpression, indent)
}
