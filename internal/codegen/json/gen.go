package json

import (
	"bytes"
	"context"
	ejson "encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

func parseOptions(req *plugin.CodeGenRequest) (*plugin.JSONCode, error) {
	if req.Settings == nil {
		return new(plugin.JSONCode), nil
	}
	if req.Settings.Codegen != nil {
		if len(req.Settings.Codegen.Options) != 0 {
			var options *plugin.JSONCode
			dec := ejson.NewDecoder(bytes.NewReader(req.Settings.Codegen.Options))
			dec.DisallowUnknownFields()
			if err := dec.Decode(&options); err != nil {
				return options, fmt.Errorf("unmarshalling options: %s", err)
			}
			return options, nil
		}
	}
	if req.Settings.Json != nil {
		return req.Settings.Json, nil
	}
	return new(plugin.JSONCode), nil
}

func Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
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
	var rm ejson.RawMessage = data
	blob, err := ejson.MarshalIndent(rm, "", indent)
	if err != nil {
		return nil, err
	}
	return &plugin.CodeGenResponse{
		Files: []*plugin.File{
			{
				Name:     filename,
				Contents: append(blob, '\n'),
			},
		},
	}, nil
}
