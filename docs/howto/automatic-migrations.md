# Automatic migration planning

[Atlas](https://atlasgo.io) is an open-source schema migration tool that can be used to automate schema migrations with `sqlc`. 

Atlas and `sqlc` complement each other by treating the same SQL schema file as the source of truth for what the 
database schema should look like and generating the necessary code to make that happen. Sqlc generates the 
database access code while Atlas maintains the actual database schema for you.

Atlas supports two kinds of workflows:

* Declarative migrations - use a Terraform-like `atlas schema apply --env sqlc` to apply your schema to the database.
* Automatic migration planning - use `atlas migrate diff --env sqlc` to automatically plan a migration from the 
  previous database schema to the current one.

## Getting started

Install Atlas from macOS or Linux by running:

  ```
  curl -sSf https://atlasgo.sh | sh
  ```
  
  See [this document](https://atlasgo.io/getting-started#installation) for more installation options.

Next, create a file named "atlas.hcl" to configure Atlas:

```hcl
env "sqlc" {
  src = "file://schema.sql"
  dev = "docker://mysql/8/dev"
  // Postgres:  dev = "docker://postgres/15/dev?search_path=public"
  // SQLite:    dev = "sqlite://dev?mode=memory"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
```

### Option 1: Apply the schema to a target database

```
atlas schema apply --env sqlc -u "mysql://root:password@localhost:3306/mydb"
```

See the docs for more [URL format examples](https://atlasgo.io/concepts/url) and information
about [declarative migrations](https://atlasgo.io/versioned/apply).

### Option 2: Use Atlas to automatically plan the next migration for you

```
atlas migrate diff --env sqlc 
```

See the docs for [automated migration planning](https://atlasgo.io/versioned/diff).

## Other migration tools

Atlas can [execute migrations](https://atlasgo.io/versioned/apply) itself or just act as
an automatic migration planner for other migration tools like [Flyway](https://flywaydb.org/),
[Liquibase](https://www.liquibase.org/), or [golang-migrate](https://github.com/golang-migrate/migrate).

To change the output format of the migration files, edit the `migration` block in `atlas.hcl`:

```diff
env "sqlc" {
  src = "file://schema.sql"
  dev = "docker://mysql/8/dev"
  // Postgres:  dev = "docker://postgres/15/dev?search_path=public"
  // SQLite:    dev = "sqlite://dev?mode=memory"
  migration {
    dir = "file://migrations"
+   format = golang-migrate // or flyway, liquibase, goose, etc.
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
```
See [this guide](https://atlasgo.io/guides/migration-tools/golang-migrate) for an example of how to use 
Atlas to plan migrations for golang-migrate.
