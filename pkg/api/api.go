package api

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sqlc-dev/sqlc/internal/cmd"
)

type Options struct {
	Dir      string // working directory for relative paths
	Filename string
	Options  *cmd.Options
}

type Diagnostic struct {
	File     string
	Line     int
	Column   int
	Severity string // "error" | "warning" | "info"
	Message  string
	Raw      string
}

type Report struct {
	Stdout      string
	Stderr      string
	Diagnostics []Diagnostic
}

func Generate(ctx context.Context, opt Options) (map[string]string, error) {
	return cmd.Generate(ctx, opt.Dir, opt.Filename, opt.Options)
}

func Verify(ctx context.Context, opt Options) error {
	return cmd.Verify(ctx, opt.Dir, opt.Filename, opt.Options)
}
