package poet

import (
	"go/format"
	"strings"
)

// Render converts a File to formatted Go source code.
func Render(f *File) ([]byte, error) {
	var b strings.Builder
	renderFile(&b, f)
	return format.Source([]byte(b.String()))
}

func renderFile(b *strings.Builder, f *File) {
	// Build tags
	if f.BuildTags != "" {
		b.WriteString("//go:build ")
		b.WriteString(f.BuildTags)
		b.WriteString("\n\n")
	}

	// File comments
	for _, c := range f.Comments {
		b.WriteString(c)
		b.WriteString("\n")
	}

	// Package
	if len(f.Comments) > 0 {
		b.WriteString("\n")
	}
	b.WriteString("package ")
	b.WriteString(f.Package)
	b.WriteString("\n")

	// Imports
	hasImports := false
	for _, group := range f.ImportGroups {
		if len(group) > 0 {
			hasImports = true
			break
		}
	}
	if hasImports {
		b.WriteString("\nimport (\n")
		first := true
		for _, group := range f.ImportGroups {
			if len(group) == 0 {
				continue
			}
			if !first {
				b.WriteString("\n")
			}
			first = false
			for _, imp := range group {
				b.WriteString("\t")
				if imp.Alias != "" {
					b.WriteString(imp.Alias)
					b.WriteString(" ")
				}
				b.WriteString("\"")
				b.WriteString(imp.Path)
				b.WriteString("\"\n")
			}
		}
		b.WriteString(")\n")
	}

	// Declarations
	for _, d := range f.Decls {
		b.WriteString("\n")
		renderDecl(b, d)
	}
}

func renderDecl(b *strings.Builder, d Decl) {
	switch d := d.(type) {
	case Raw:
		b.WriteString(d.Code)
	case Const:
		renderConst(b, d, "")
	case ConstBlock:
		renderConstBlock(b, d)
	case Var:
		renderVar(b, d, "")
	case VarBlock:
		renderVarBlock(b, d)
	case TypeDef:
		renderTypeDef(b, d)
	case Func:
		renderFunc(b, d)
	}
}

func renderConst(b *strings.Builder, c Const, indent string) {
	if c.Comment != "" {
		writeComment(b, c.Comment, indent)
	}
	b.WriteString(indent)
	if indent == "" {
		b.WriteString("const ")
	}
	b.WriteString(c.Name)
	if c.Type != "" {
		b.WriteString(" ")
		b.WriteString(c.Type)
	}
	if c.Value != "" {
		b.WriteString(" = ")
		b.WriteString(c.Value)
	}
	b.WriteString("\n")
}

func renderConstBlock(b *strings.Builder, cb ConstBlock) {
	b.WriteString("const (\n")
	for _, c := range cb.Consts {
		renderConst(b, c, "\t")
	}
	b.WriteString(")\n")
}

func renderVar(b *strings.Builder, v Var, indent string) {
	if v.Comment != "" {
		writeComment(b, v.Comment, indent)
	}
	b.WriteString(indent)
	if indent == "" {
		b.WriteString("var ")
	}
	b.WriteString(v.Name)
	if v.Type != "" {
		b.WriteString(" ")
		b.WriteString(v.Type)
	}
	if v.Value != "" {
		b.WriteString(" = ")
		b.WriteString(v.Value)
	}
	b.WriteString("\n")
}

func renderVarBlock(b *strings.Builder, vb VarBlock) {
	b.WriteString("var (\n")
	for _, v := range vb.Vars {
		renderVar(b, v, "\t")
	}
	b.WriteString(")\n")
}

func renderTypeDef(b *strings.Builder, t TypeDef) {
	if t.Comment != "" {
		writeComment(b, t.Comment, "")
	}
	b.WriteString("type ")
	b.WriteString(t.Name)
	b.WriteString(" ")
	renderTypeExpr(b, t.Type)
	b.WriteString("\n")
}

func renderTypeExpr(b *strings.Builder, t TypeExpr) {
	switch t := t.(type) {
	case Struct:
		renderStruct(b, t)
	case Interface:
		renderInterface(b, t)
	case TypeName:
		b.WriteString(t.Name)
	}
}

func renderStruct(b *strings.Builder, s Struct) {
	b.WriteString("struct {\n")
	for _, f := range s.Fields {
		if f.Comment != "" {
			writeComment(b, f.Comment, "\t")
		}
		b.WriteString("\t")
		b.WriteString(f.Name)
		b.WriteString(" ")
		b.WriteString(f.Type)
		if f.Tag != "" {
			b.WriteString(" `")
			b.WriteString(f.Tag)
			b.WriteString("`")
		}
		if f.TrailingComment != "" {
			b.WriteString(" // ")
			b.WriteString(f.TrailingComment)
		}
		b.WriteString("\n")
	}
	b.WriteString("}")
}

func renderInterface(b *strings.Builder, iface Interface) {
	b.WriteString("interface {\n")
	for _, m := range iface.Methods {
		if m.Comment != "" {
			writeComment(b, m.Comment, "\t")
		}
		b.WriteString("\t")
		b.WriteString(m.Name)
		b.WriteString("(")
		renderParams(b, m.Params)
		b.WriteString(")")
		if len(m.Results) > 0 {
			b.WriteString(" ")
			if len(m.Results) == 1 && m.Results[0].Name == "" {
				b.WriteString(m.Results[0].Type)
			} else {
				b.WriteString("(")
				renderParams(b, m.Results)
				b.WriteString(")")
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("}")
}

func renderFunc(b *strings.Builder, f Func) {
	if f.Comment != "" {
		writeComment(b, f.Comment, "")
	}
	b.WriteString("func ")
	if f.Recv != nil {
		b.WriteString("(")
		b.WriteString(f.Recv.Name)
		b.WriteString(" ")
		if f.Recv.Pointer {
			b.WriteString("*")
		}
		b.WriteString(f.Recv.Type)
		b.WriteString(") ")
	}
	b.WriteString(f.Name)
	b.WriteString("(")
	renderParams(b, f.Params)
	b.WriteString(")")
	if len(f.Results) > 0 {
		b.WriteString(" ")
		if len(f.Results) == 1 && f.Results[0].Name == "" {
			if f.Results[0].Pointer {
				b.WriteString("*")
			}
			b.WriteString(f.Results[0].Type)
		} else {
			b.WriteString("(")
			renderParams(b, f.Results)
			b.WriteString(")")
		}
	}
	b.WriteString(" {\n")
	renderStmts(b, f.Stmts, "\t")
	b.WriteString("}\n")
}

func renderParams(b *strings.Builder, params []Param) {
	for i, p := range params {
		if i > 0 {
			b.WriteString(", ")
		}
		if p.Name != "" {
			b.WriteString(p.Name)
			b.WriteString(" ")
		}
		if p.Pointer {
			b.WriteString("*")
		}
		b.WriteString(p.Type)
	}
}

func writeComment(b *strings.Builder, comment, indent string) {
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		b.WriteString(indent)
		// If line already starts with //, write as-is
		if strings.HasPrefix(line, "//") {
			b.WriteString(line)
		} else {
			b.WriteString("// ")
			b.WriteString(line)
		}
		b.WriteString("\n")
	}
}

func renderStmts(b *strings.Builder, stmts []Stmt, indent string) {
	for _, s := range stmts {
		renderStmt(b, s, indent)
	}
}

func renderStmt(b *strings.Builder, s Stmt, indent string) {
	switch s := s.(type) {
	case RawStmt:
		b.WriteString(s.Code)
	case Return:
		renderReturn(b, s, indent)
	case For:
		renderFor(b, s, indent)
	case If:
		renderIf(b, s, indent)
	case Switch:
		renderSwitch(b, s, indent)
	}
}

func renderReturn(b *strings.Builder, r Return, indent string) {
	b.WriteString(indent)
	b.WriteString("return")
	if len(r.Values) > 0 {
		b.WriteString(" ")
		for i, v := range r.Values {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(v)
		}
	}
	b.WriteString("\n")
}

func renderFor(b *strings.Builder, f For, indent string) {
	b.WriteString(indent)
	b.WriteString("for ")
	if f.Range != "" {
		b.WriteString(f.Range)
	} else {
		if f.Init != "" {
			b.WriteString(f.Init)
		}
		b.WriteString("; ")
		b.WriteString(f.Cond)
		b.WriteString("; ")
		if f.Post != "" {
			b.WriteString(f.Post)
		}
	}
	b.WriteString(" {\n")
	renderStmts(b, f.Body, indent+"\t")
	b.WriteString(indent)
	b.WriteString("}\n")
}

func renderIf(b *strings.Builder, i If, indent string) {
	b.WriteString(indent)
	b.WriteString("if ")
	if i.Init != "" {
		b.WriteString(i.Init)
		b.WriteString("; ")
	}
	b.WriteString(i.Cond)
	b.WriteString(" {\n")
	renderStmts(b, i.Body, indent+"\t")
	b.WriteString(indent)
	b.WriteString("}")
	if len(i.Else) > 0 {
		b.WriteString(" else {\n")
		renderStmts(b, i.Else, indent+"\t")
		b.WriteString(indent)
		b.WriteString("}")
	}
	b.WriteString("\n")
}

func renderSwitch(b *strings.Builder, s Switch, indent string) {
	b.WriteString(indent)
	b.WriteString("switch ")
	if s.Init != "" {
		b.WriteString(s.Init)
		b.WriteString("; ")
	}
	b.WriteString(s.Expr)
	b.WriteString(" {\n")
	for _, c := range s.Cases {
		b.WriteString(indent)
		if len(c.Values) == 0 {
			b.WriteString("default:\n")
		} else {
			b.WriteString("case ")
			if len(c.Values) == 1 {
				b.WriteString(c.Values[0])
			} else {
				// Multiple values: put each on its own line
				for i, v := range c.Values {
					if i > 0 {
						b.WriteString(",\n")
						b.WriteString(indent)
						b.WriteString("\t")
					}
					b.WriteString(v)
				}
			}
			b.WriteString(":\n")
		}
		renderStmts(b, c.Body, indent+"\t")
	}
	b.WriteString(indent)
	b.WriteString("}\n")
}
