package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sqlc-dev/sqlc/internal/codegen/json"
	"github.com/sqlc-dev/sqlc/internal/plugin"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error generating JSON: %s", err)
		os.Exit(2)
	}
}

func run() error {
	var req plugin.GenerateRequest
	reqBlob, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(reqBlob, &req); err != nil {
		return err
	}
	resp, err := json.Generate(context.Background(), &req)
	if err != nil {
		return err
	}
	respBlob, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(os.Stdout)
	if _, err := w.Write(respBlob); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
