package ydb

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	parser "github.com/ydb-platform/yql-parsers/go"
)

type objectRefProvider interface {
	antlr.ParserRuleContext
	Object_ref() parser.IObject_refContext
}



func parseTableName(ctx objectRefProvider) *ast.TableName {
	return parseObjectRef(ctx.Object_ref())
}

func parseObjectRef(r parser.IObject_refContext) *ast.TableName {
	if r == nil {
		return nil
	}
	ref := r.(*parser.Object_refContext)

	parts := []string{}

	if cl := ref.Cluster_expr(); cl != nil {
		parts = append(parts, parseClusterExpr(cl))
	}

	if idOrAt := ref.Id_or_at(); idOrAt != nil {
		parts = append(parts, parseIdOrAt(idOrAt))
	}

	objectName := strings.Join(parts, ".")

	return &ast.TableName{
		Schema: "",
		Name:   identifier(objectName),
	}
}

func parseClusterExpr(ctx parser.ICluster_exprContext) string {
	if ctx == nil {
		return ""
	}
	return identifier(ctx.GetText())
}

func parseIdOrAt(ctx parser.IId_or_atContext) string {
	if ctx == nil {
		return ""
	}
	idOrAt := ctx.(*parser.Id_or_atContext)

	if ao := idOrAt.An_id_or_type(); ao != nil {
		return identifier(parseAnIdOrType(ao))
	}
	return ""
}

func parseAnIdOrType(ctx parser.IAn_id_or_typeContext) string {
	if ctx == nil {
		return ""
	}
	anId := ctx.(*parser.An_id_or_typeContext)

	if anId.Id_or_type() != nil {
		return identifier(parseIdOrType(anId.Id_or_type()))
	}

	if anId.STRING_VALUE() != nil {
		return identifier(anId.STRING_VALUE().GetText())
	}

	return ""
}

func parseIdOrType(ctx parser.IId_or_typeContext) string {
	if ctx == nil {
		return ""
	}
	Id := ctx.(*parser.Id_or_typeContext)
	if Id.Id() != nil {
		return identifier(parseIdTable(Id.Id()))
	}

	return ""
}

func parseAnId(ctx parser.IAn_idContext) string {
	if id := ctx.Id(); id != nil {
		return id.GetText()
	} else if str := ctx.STRING_VALUE(); str != nil {
		return str.GetText()
	}
	return ""
}

func parseAnIdSchema(ctx parser.IAn_id_schemaContext) string {
	if ctx == nil {
		return ""
	}
	if id := ctx.Id_schema(); id != nil {
		return id.GetText()
	} else if str := ctx.STRING_VALUE(); str != nil {
		return str.GetText()
	}
	return ""
}

func parseIdTable(ctx parser.IIdContext) string {
	if ctx == nil {
		return ""
	}
	return ctx.GetText()
}

func parseIntegerValue(text string) (int64, error) {
	text = strings.ToLower(text)
	base := 10

	switch {
	case strings.HasPrefix(text, "0x"):
		base = 16
		text = strings.TrimPrefix(text, "0x")

	case strings.HasPrefix(text, "0o"):
		base = 8
		text = strings.TrimPrefix(text, "0o")

	case strings.HasPrefix(text, "0b"):
		base = 2
		text = strings.TrimPrefix(text, "0b")
	}

	// debug!!!
	text = strings.TrimRight(text, "pulstibn")

	return strconv.ParseInt(text, base, 64)
}
