# CLI

```sh
Usage:
  sqlc [command]

Available Commands:
  compile     Statically check SQL for syntax and type errors
  completion  Generate the autocompletion script for the specified shell
  createdb    Create an ephemeral database
  diff        Compare the generated files to the existing files
  generate    Generate source code from SQL
  help        Help about any command
  init        Create an empty sqlc.yaml settings file
  push        Push the schema, queries, and configuration for this project
  verify      Verify schema, queries, and configuration for this project
  version     Print the sqlc version number
  vet         Vet examines queries

Flags:
  -f, --file string    specify an alternate config file (default: sqlc.yaml)
  -h, --help           help for sqlc
      --no-database    disable database connections (default: false)
      --no-remote      disable remote execution (default: false)

Use "sqlc [command] --help" for more information about a command.
```
