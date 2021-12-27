package poet

import "github.com/kyleconroy/sqlc/internal/kotlin/ast"


func Comment(text string) *ast.Node {
	return &ast.Node{
		Node: &ast.Node_Comment{
			Comment: &ast.Comment{
				Text: text,
			},
		},
	}
}

func DotQualifiedExpression(rec, sel *ast.Node) *ast.Node {
	return &ast.Node{
		Node: &ast.Node_DotQualifiedExpression{
			DotQualifiedExpression: &ast.DotQualifiedExpression{
				Receiver: rec,
				Selector: sel,
			},
		},
	}
}

func NameReferenceExpression(name string) *ast.Node {
	return &ast.Node{
		Node: &ast.Node_NameReferenceExpression{
			NameReferenceExpression: &ast.NameReferenceExpression{
				Name: name,
			},
		},
	}
}

func PackageDirective(name *ast.Node) *ast.Node {
	return &ast.Node{
		Node: &ast.Node_PackageDirective{
			PackageDirective: &ast.PackageDirective{
				Name: name,
			},
		},
	}
}

