package json

import (
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

func Generate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	indent := ""
	if req.Settings != nil && req.Settings.Json != nil {
		indent = req.Settings.Json.Indent
	}
	m := protojson.MarshalOptions{
		Indent: indent,
	}
	blob, err := m.Marshal(req)
	if err != nil {
		return nil, err
	}
	return &plugin.CodeGenResponse{
		Files: []*plugin.File{
			{
				Name:     "codegen_request.json",
				Contents: blob,
			},
		},
	}, nil
}
