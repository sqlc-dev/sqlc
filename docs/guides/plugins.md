# Authoring plugins

## Process plugins

To use process-based plugins, you must be using [Version
2](../reference/config.html) of the configuration file. The top-level `plugins`
array defines the available plugins and maps them to an executable on the system.

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

### Security

Process-based plugins offer minimal security. Only use plugins that you trust.
Better yet, only use plugins that you've written yourself.

## WASM plugins

*Coming soon!*
