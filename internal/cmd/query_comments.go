package cmd

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

const (
	queryCommentFormatMarginalia   = "marginalia"
	queryCommentFormatSQLCommenter = "sqlcommenter"
)

func applyQueryComments(req *plugin.GenerateRequest, opts config.QueryComments) {
	if !opts.Enabled {
		return
	}
	for _, query := range req.Queries {
		if query.Text == "" {
			continue
		}
		comment := queryComment(query, opts)
		if comment == "" {
			continue
		}
		query.Text = comment + " " + query.Text
	}
}

func queryComment(query *plugin.Query, opts config.QueryComments) string {
	tags := opts.Tags
	if len(tags) == 0 {
		tags = []string{"name"}
	}

	parts := make([]string, 0, len(tags))
	for _, tag := range tags {
		value := queryCommentValue(query, tag)
		if value == "" {
			continue
		}
		key := "sqlc_" + tag
		if opts.Format == queryCommentFormatMarginalia {
			parts = append(parts, key+":"+escapeQueryCommentValue(value))
		} else {
			parts = append(parts, key+"='"+escapeQueryCommentValue(value)+"'")
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "/*" + strings.Join(parts, ",") + "*/"
}

func queryCommentValue(query *plugin.Query, tag string) string {
	switch tag {
	case "name":
		return query.Name
	case "cmd":
		return query.Cmd
	case "filename":
		return query.Filename
	default:
		return ""
	}
}

func escapeQueryCommentValue(value string) string {
	value = strings.ReplaceAll(value, "*/", "* /")
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "'", "%27")
	value = strings.ReplaceAll(value, ",", "%2C")
	value = strings.ReplaceAll(value, ":", "%3A")
	return value
}
