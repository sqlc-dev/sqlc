# Migrating off hosted managed databases
 
Starting in sqlc 1.27.0, [managed databases](../docs/managed-databases.md) will require a database server URI in the configuration file.

This guide walks you through migrating to a locally running database server.

## Run a database server locally

There are many options for running a database server locally, but this guide
will use [Docker Compose](https://docs.docker.com/compose/), as it can support
both MySQL and PostgreSQL.

If you're using macOS and PostgreSQL, [Postgres.app](https://postgresapp.com/) is also a good option.

For MySQL, create a `docker-compose.yml` file with the following contents:

```yaml
version: "3.8"
services:
  mysql:
    image: "mysql/mysql-server:8.0"
    ports:
      - "3306:3306"
    restart: always
    environment:
      MYSQL_DATABASE: dinotest
      MYSQL_ROOT_PASSWORD: mysecretpassword
      MYSQL_ROOT_HOST: '%'
```

For PostgreSQL, create a `docker-compose.yml` file with the following contents:

```yaml
version: "3.8"
services:
  postgresql:
    image: "postgres:16"
    ports:
      - "5432:5432"
    restart: always
    environment:
      POSTGRES_DB: postgres
      POSTGRES_PASSWORD: mysecretpassword
      POSTGRES_USER: postgres
```

```sh
docker compose up -d
```

## Upgrade sqlc

You must be running sqlc v1.30.0 or greater to have access to the `servers`
configuration.

## Add servers to configuration

```diff
version: '2'
cloud:
  project: '<PROJECT_ID>'
+ servers:
+ - name: mysql
+   uri: mysql://localhost:3306
+ - name: postgres
+   uri: postgres://localhost:5432/postgres?sslmode=disable
```

## Re-generate the code

Run `sqlc generate`. A database with the `sqlc_managed_` prefix will be automatically created and used for query analysis. 
