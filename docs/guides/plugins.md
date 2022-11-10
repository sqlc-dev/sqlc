# Authoring plugins

To use plugins, you must be using [Version 2](../reference/config.html) of
the configuration file. The top-level `plugins` array defines the available
plugins.

## WASM plugins

> WASM plugins are fully sandboxed. Plugins do not have access to the network,
> filesystem, or environment variables.

In the `codegen` section, the `out` field dictates what directory will contain
the new files. The `plugin` key must reference a plugin defined in the
top-level `plugins` map. The `options` are serialized to a string and passed on
to the plugin itself.


```json
{
  "version": "2",
  "plugins": [
    {
      "name": "greeter",
      "wasm": {
        "url": "https://github.com/kyleconroy/sqlc-gen-greeter/releases/download/v0.1.0/sqlc-gen-greeter.wasm",
        "sha256": "afc486dac2068d741d7a4110146559d12a013fd0286f42a2fc7dcd802424ad07"
      }
    }
  ],
  "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      "codegen": [
        {
          "out": "gen",
          "plugin": "greeter"
        }
      ]
    }
  ]
}
```

For a complete working example see the following files:
- [sqlc-gen-greeter](https://github.com/kyleconroy/sqlc-gen-greeter)
  - A WASM plugin (written in Rust) that outputs a friendly message
- [wasm_plugin_sqlc_gen_greeter](https://github.com/kyleconroy/sqlc/tree/main/internal/endtoend/testdata/wasm_plugin_sqlc_gen_greeter)
  - An example project showing how to use a WASM plugin

## Process plugins

> Process-based plugins offer minimal security. Only use plugins that you
> trust. Better yet, only use plugins that you've written yourself.

In the `codegen` section, the `out` field dictates what directory will contain
the new files. The `plugin` key must reference a plugin defined in the
top-level `plugins` map. The `options` are serialized to a string and passed on
to the plugin itself.

```json
{
  "version": "2",
  "plugins": [
    {
      "name": "jsonb",
      "process": {
        "cmd": "sqlc-gen-json"
      }
    }
  ],
   "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      "codegen": [
        {
          "out": "gen",
          "plugin": "jsonb",
          "options": {
            "indent": "  ",
            "filename": "codegen.json"
          }
        }
      ]
    }
  ]
}
```

For a complete working example see the following files:
- [sqlc-gen-json](https://github.com/kyleconroy/sqlc/tree/main/cmd/sqlc-gen-json)
  - A process-based plugin that serializes the CodeGenRequest to JSON
- [process_plugin_sqlc_gen_json](https://github.com/kyleconroy/sqlc/tree/main/internal/endtoend/testdata/process_plugin_sqlc_gen_json)
  - An example project showing how to use a process-based plugin
