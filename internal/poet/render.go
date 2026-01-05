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
			b.WriteString(f.Results[0].Type)
		} else {
			b.WriteString("(")
			renderParams(b, f.Results)
			b.WriteString(")")
		}
	}
	b.WriteString(" {\n")
	b.WriteString(f.Body)
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
