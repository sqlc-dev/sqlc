# Migrating to sqlc-gen-kotlin
 
Starting in sqlc 1.16.0, built-in Kotlin support has been deprecated. It will
be fully removed in 1.17.0 in favor of sqlc-gen-kotlin.

This guide will walk you through migrating to the [sqlc-gen-kotlin] plugin,
which involves three steps.

1. Add the sqlc-gen-kotlin plugin
2. Migrate each package
3. Re-generate the code

## Add the sqlc-gen-kotlin plugin

In your configuration file, add a `plugins` array if you don't have one
already. Add the following configuration for the plugin:

```json
{
  "version": "2",
  "plugins": [
    {
      "name": "py",
      "wasm": {
        "url": "https://github.com/tabbed/sqlc-gen-kotlin/releases/download/v0.16.0-alpha/sqlc-gen-python.wasm",
        "sha256": "4fb54ee7d25b4d909b59a8271ebee60ad76ff17b10d61632a5ca5651e4bfe438"
      }
    }
  ]
}
```

```yaml
version: "2"
plugins:
  name: py,
  wasm:
    url: "https://github.com/tabbed/sqlc-gen-kotlin/releases/download/v0.16.0-alpha/sqlc-gen-python.wasm"
    sha256: "4fb54ee7d25b4d909b59a8271ebee60ad76ff17b10d61632a5ca5651e4bfe438"
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
        "kotlin": {
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
      kotlin:
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
