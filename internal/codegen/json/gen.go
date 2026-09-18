package json

import (
	"bytes"
	"context"
	ejson "encoding/json"
	"fmt"
	"regexp"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func parseOptions(req *plugin.GenerateRequest) (*opts, error) {
	if len(req.PluginOptions) == 0 {
		return new(opts), nil
	}

	var options *opts
	dec := ejson.NewDecoder(bytes.NewReader(req.PluginOptions))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&options); err != nil {
		return options, fmt.Errorf("unmarshalling options: %s", err)
	}
	return options, nil
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := parseOptions(req)
	if err != nil {
		return nil, err
	}

	indent := "  "
	if options.Indent != "" {
		indent = options.Indent
	}

	filename := "codegen_request.json"
	if options.Filename != "" {
		filename = options.Filename
	}

	// The output of protojson has randomized whitespace
	// https://github.com/golang/protobuf/issues/1082
	m := &protojson.MarshalOptions{
		EmitUnpopulated: true,
		Indent:          "",
		UseProtoNames:   true,
	}
	data, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = dropUnsetTypeExpr(data)
	var rm ejson.RawMessage = data
	blob, err := ejson.MarshalIndent(rm, "", indent)
	if err != nil {
		return nil, err
	}
	return &plugin.GenerateResponse{
		Files: []*plugin.File{
			{
				Name:     filename,
				Contents: append(blob, '\n'),
			},
		},
	}, nil
}

// unsetTypeExpr matches a column's type_expr member when it is unset, with
// the comma that joins it to the member before it. It is the last member
// of a column, so protojson writes it after a comma, with whatever
// whitespace it chose around the separators.
var unsetTypeExpr = regexp.MustCompile(`,\s*"type_expr"\s*:\s*null`)

// dropUnsetTypeExpr removes the type_expr member from columns that have
// none. The member is only set when the analysis core typed the column,
// and a column the legacy path typed is written the way it was before the
// member existed rather than with an explicit null. A JSON string cannot
// hold an unescaped quote, so the match never lands inside one.
func dropUnsetTypeExpr(data []byte) []byte {
	return unsetTypeExpr.ReplaceAll(data, nil)
}
