# Migrating to sqlc-gen-kotlin
 
Starting in sqlc 1.16.0, built-in Kotlin support has been deprecated. It will
be fully removed in 1.17.0 in favor of sqlc-gen-kotlin.

This guide will walk you through migrating to the [sqlc-gen-kotlin](https://github.com/sqlc-dev/sqlc-gen-kotlin) plugin,
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
      "name": "kt",
      "wasm": {
        "url": "https://downloads.sqlc.dev/plugin/sqlc-gen-kotlin_1.0.0.wasm",
        "sha256": "7620dc5d462de41fdc90e2011232c842117b416c98fd5c163d27c5738431a45c"
      }
    }
  ]
}
```

```yaml
version: "2"
plugins:
  name: "kt"
  wasm:
    url: "https://downloads.sqlc.dev/plugin/sqlc-gen-kotlin_1.0.0.wasm"
    sha256: "7620dc5d462de41fdc90e2011232c842117b416c98fd5c163d27c5738431a45c"
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
          "out": "src/main/kotlin/com/example/foo",
          "package": "com.example.foo"
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
        out: "src/main/kotlin/com/example/foo"
        package: "com.example.foo"
```

To use the plugin, you'll need to replace the `gen` mapping with the `codegen`
collection. Add the `plugin` field, setting it to `kt`. All fields other than
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
          "out": "src/main/kotlin/com/example/foo",
          "plugin": "kt",
          "options": {
            "package": "com.example.foo"
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
    - plugin: "kt"
      out: "src/main/kotlin/com/example/foo"
      options:
        package: "com.example.foo"
```

## Re-generate the code

Run `sqlc generate`. The plugin will produce the same output, so you shouldn't
see any changes. The first time `sqlc generate` is run, the plugin must be
downloaded and compiled, resulting in a slightly longer runtime. Subsequent
`generate` calls will be fast.
