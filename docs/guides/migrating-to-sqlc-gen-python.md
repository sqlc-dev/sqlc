# Migrating to sqlc-gen-python
 
Starting in sqlc 1.16.0, built-in Python support has been deprecated. It will
be fully removed in 1.17.0 in favor of sqlc-gen-python.

This guide will walk you through migrating to the [sqlc-gen-python](https://github.com/sqlc-dev/sqlc-gen-python) plugin,
which involves three steps.

1. Add the sqlc-gen-python plugin
2. Migrate each package
3. Re-generate the code

## Add the sqlc-gen-python plugin

In your configuration file, add a `plugins` array if you don't have one
already. Add the following configuration for the plugin:

```json
{
  "version": "2",
  "plugins": [
    {
      "name": "py",
      "wasm": {
        "url": "https://downloads.sqlc.dev/plugin/sqlc-gen-python_1.0.0.wasm",
        "sha256": "aca83e1f59f8ffdc604774c2f6f9eb321a2b23e07dc83fc12289d25305fa065b"
      }
    }
  ]
}
```

```yaml
version: "2"
plugins:
  - name: "py"
    wasm:
      url: "https://downloads.sqlc.dev/plugin/sqlc-gen-python_1.0.0.wasm"
      sha256: "aca83e1f59f8ffdc604774c2f6f9eb321a2b23e07dc83fc12289d25305fa065b"
```

## Migrate each package

Your package configuration should currently looks something like this for JSON.

```json
  "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      "gen": {
        "python": {
          "out": "src",
          "package": "foo",
          "emit_sync_querier": true,
          "emit_async_querier": true,
          "query_parameter_limit": 5
        }
      }
    }
  ]
```

Or this if you're using YAML.

```yaml
  sql:
  - schema: "schema.sql"
    queries: "query.sql"
    engine: "postgresql"
    gen:
      python:
        out: "src"
        package: "foo"
        emit_sync_querier: true
        emit_async_querier: true
        query_parameter_limit: 5
```

To use the plugin, you'll need to replace the `gen` mapping with the `codegen`
collection. Add the `plugin` field, setting it to `py`. All fields other than
`out` need to be moved into the `options` mapping.

After you're done, it should look like this for JSON.

```json
  "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      "codegen": [
        {
          "out": "src",
          "plugin": "py",
          "options": {
            "package": "authors",
            "emit_sync_querier": true,
            "emit_async_querier": true,
            "query_parameter_limit": 5
          }
        }
      ]
    }
  ]
```

Or this for YAML.

```yaml
  sql:
  - schema: "schema.sql"
    queries: "query.sql"
    engine: "postgresql"
    codegen:
    - plugin: "py"
      out: "src"
      options:
        package: "foo"
        emit_sync_querier: true
        emit_async_querier: true
        query_parameter_limit: 5
```

## Re-generate the code

Run `sqlc generate`. The plugin will produce the same output, so you shouldn't
see any changes. The first time `sqlc generate` is run, the plugin must be
downloaded and compiled, resulting in a slightly longer runtime. Subsequent
`generate` calls will be fast.
