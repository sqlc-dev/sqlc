package json

import (
	ejson "encoding/json"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

func Generate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	indent := ""
	if req.Settings != nil && req.Settings.Json != nil {
		indent = req.Settings.Json.Indent
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
				Name:     "codegen_request.json",
				Contents: append(blob, '\n'),
			},
		},
	}, nil
}
