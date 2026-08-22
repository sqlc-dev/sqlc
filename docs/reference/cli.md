# CLI

```sh
Usage:
  sqlc [command]

Available Commands:
  analyze     Analyze a query against a schema and output the result columns and parameters
  compile     Statically check SQL for syntax and type errors
  completion  Generate the autocompletion script for the specified shell
  createdb    Create an ephemeral database
  diff        Compare the generated files to the existing files
  fmt         Format query files
  generate    Generate source code from SQL
  help        Help about any command
  init        Create an empty sqlc.yaml settings file
  parse       Parse SQL and output the AST as JSON
  push        Push the schema, queries, and configuration for this project
  verify      Verify schema, queries, and configuration for this project
  version     Print the sqlc version number
  vet         Vet examines queries

Flags:
  -f, --file string    specify an alternate config file (default: sqlc.yaml)
  -h, --help           help for sqlc
      --no-database    disable database connections (default: false)

Use "sqlc [command] --help" for more information about a command.
```
