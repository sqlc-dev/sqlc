package main

import (
	"bufio"
	"io"
	"os"

	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func main() {
	var req plugin.CodeGenRequest
	reqBlob, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(10)
	}
	if err := req.UnmarshalVT(reqBlob); err != nil {
		os.Exit(11)
	}
	resp, err := golang.Generate(&req)
	if err != nil {
		os.Exit(12)
	}
	respBlob, err := resp.MarshalVT()
	if err != nil {
		os.Exit(13)
	}
	w := bufio.NewWriter(os.Stdout)
	if _, err := w.Write(respBlob); err != nil {
		os.Exit(14)
	}
	if err := w.Flush(); err != nil {
		os.Exit(15)
	}
}
