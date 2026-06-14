package engine

import (
	"bytes"
	"io"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestRun_acceptsFullGRPCMethodName(t *testing.T) {
	called := false
	h := Handler{
		Parse: func(req *ParseRequest) (*ParseResponse, error) {
			called = true
			return &ParseResponse{}, nil
		},
	}
	in, err := proto.Marshal(&ParseRequest{Sql: "SELECT 1"})
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err = run(h, []string{"plugin", EngineService_Parse_FullMethodName}, bytes.NewReader(in), &stdout, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("Parse not invoked for full gRPC method argv")
	}
	var resp ParseResponse
	if err := proto.Unmarshal(stdout.Bytes(), &resp); err != nil {
		t.Fatalf("stdout protobuf: %v", err)
	}
}

func TestRun_acceptsLegacyParseArgv(t *testing.T) {
	called := false
	h := Handler{
		Parse: func(req *ParseRequest) (*ParseResponse, error) {
			called = true
			return &ParseResponse{}, nil
		},
	}
	in, err := proto.Marshal(&ParseRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err = run(h, []string{"plugin", "parse"}, bytes.NewReader(in), &stdout, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("Parse not invoked for legacy parse argv")
	}
}
