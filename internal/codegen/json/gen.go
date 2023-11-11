package json

import (
	"bytes"
	"context"
	ejson "encoding/json"
	"fmt"

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
