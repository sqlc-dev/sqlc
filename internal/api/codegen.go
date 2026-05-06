package api

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/trace"

	"google.golang.org/grpc"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang"
	genjson "github.com/sqlc-dev/sqlc/internal/codegen/json"
	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/config/convert"
	"github.com/sqlc-dev/sqlc/internal/ext"
	"github.com/sqlc-dev/sqlc/internal/ext/process"
	"github.com/sqlc-dev/sqlc/internal/ext/wasm"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func findPlugin(conf config.Config, name string) (*config.Plugin, error) {
	for _, plug := range conf.Plugins {
		if plug.Name == name {
			return &plug, nil
		}
	}
	return nil, fmt.Errorf("plugin not found")
}

func codegen(ctx context.Context, combo config.CombinedSettings, sql outputPair, result *compiler.Result) (string, *plugin.GenerateResponse, error) {
	defer trace.StartRegion(ctx, "codegen").End()
	req := codeGenRequest(result, combo)
	var handler grpc.ClientConnInterface
	var out string
	switch {
	case sql.Plugin != nil:
		out = sql.Plugin.Out
		plug, err := findPlugin(combo.Global, sql.Plugin.Plugin)
		if err != nil {
			return "", nil, fmt.Errorf("plugin not found: %s", err)
		}

		switch {
		case plug.Process != nil:
			handler = &process.Runner{
				Cmd:    plug.Process.Cmd,
				Env:    plug.Env,
				Format: plug.Process.Format,
			}
		case plug.WASM != nil:
			handler = &wasm.Runner{
				URL:    plug.WASM.URL,
				SHA256: plug.WASM.SHA256,
				Env:    plug.Env,
			}
		default:
			return "", nil, fmt.Errorf("unsupported plugin type")
		}

		opts, err := convert.YAMLtoJSON(sql.Plugin.Options)
		if err != nil {
			return "", nil, fmt.Errorf("invalid plugin options: %w", err)
		}
		req.PluginOptions = opts

		global, found := combo.Global.Options[plug.Name]
		if found {
			opts, err := convert.YAMLtoJSON(global)
			if err != nil {
				return "", nil, fmt.Errorf("invalid global options: %w", err)
			}
			req.GlobalOptions = opts
		}

	case sql.Gen.Go != nil:
		out = combo.Go.Out
		handler = ext.HandleFunc(golang.Generate)
		opts, err := json.Marshal(sql.Gen.Go)
		if err != nil {
			return "", nil, fmt.Errorf("opts marshal failed: %w", err)
		}
		req.PluginOptions = opts

		if combo.Global.Overrides.Go != nil {
			opts, err := json.Marshal(combo.Global.Overrides.Go)
			if err != nil {
				return "", nil, fmt.Errorf("opts marshal failed: %w", err)
			}
			req.GlobalOptions = opts
		}

	case sql.Gen.JSON != nil:
		out = combo.JSON.Out
		handler = ext.HandleFunc(genjson.Generate)
		opts, err := json.Marshal(sql.Gen.JSON)
		if err != nil {
			return "", nil, fmt.Errorf("opts marshal failed: %w", err)
		}
		req.PluginOptions = opts

	default:
		return "", nil, fmt.Errorf("missing language backend")
	}
	client := plugin.NewCodegenServiceClient(handler)
	resp, err := client.Generate(ctx, req)
	return out, resp, err
}
