# Claude Code Development Guide for sqlc

This document provides essential information for working with the sqlc codebase, including testing, development workflow, and code structure.

## Quick Start

### Prerequisites

- **Go 1.26.0+** - Required for building and testing
- **Docker & Docker Compose** - Required for integration tests with databases (local development)
- **Git** - For version control

## Database Setup with sqlc-test-setup

The `sqlc-test-setup` tool (`cmd/sqlc-test-setup/`) automates installing and starting PostgreSQL and MySQL for tests. Both commands are idempotent and safe to re-run.

### Install databases

```bash
go run ./cmd/sqlc-test-setup install
```

This will:
- Configure the apt proxy (if `http_proxy` is set, e.g. in Claude Code remote environments)
- Install PostgreSQL via apt
- Download and install MySQL 9 from Oracle's deb bundle
- Resolve all dependencies automatically
- Skip anything already installed

### Start databases

```bash
go run ./cmd/sqlc-test-setup start
```

This will:
- Start PostgreSQL and configure password auth (`postgres`/`postgres`)
- Start MySQL via `mysqld_safe` and set root password (`mysecretpassword`)
- Verify both connections
- Skip steps that are already done (running services, existing config)

Connection URIs after start:
- PostgreSQL: `postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable`
- MySQL: `root:mysecretpassword@tcp(127.0.0.1:3306)/mysql`

### Run tests

```bash
# Full test suite (requires databases running)
go test --tags=examples -timeout 20m ./...
```

## Running Tests

### Basic Unit Tests (No Database Required)

```bash
go test ./...
```

### Full Test Suite with Docker (Local Development)

```bash
docker compose up -d
go test --tags=examples -timeout 20m ./...
```

### Full Test Suite without Docker (Remote / CI)

```bash
go run ./cmd/sqlc-test-setup install
go run ./cmd/sqlc-test-setup start
go test --tags=examples -timeout 20m ./...
```

### Running Specific Tests

```bash
# Test a specific package
go test ./internal/config

# Run with verbose output
go test -v ./internal/config

# Run a specific test function
go test -v ./internal/config -run TestConfig

# Run with race detector (recommended for concurrency changes)
go test -race ./internal/config
```

## Test Types

### 1. Unit Tests

- **Location:** Throughout the codebase as `*_test.go` files
- **Run without:** Database or external dependencies
- **Examples:**
  - `/internal/config/config_test.go` - Configuration parsing
  - `/internal/compiler/selector_test.go` - Compiler logic
  - `/internal/metadata/metadata_test.go` - Query metadata parsing

### 2. End-to-End Tests

- **Location:** `/internal/endtoend/`
- **Requirements:** `--tags=examples` flag and running databases
- **Tests:**
  - `TestExamples` - Main end-to-end tests
  - `TestReplay` - Replay tests
  - `TestFormat` - Code formatting tests
  - `TestJsonSchema` - JSON schema validation
  - `TestExamplesVet` - Static analysis tests

### 3. Example Tests

- **Location:** `/examples/` directory
- **Requirements:** Tagged with "examples", requires live databases
- **Databases:** PostgreSQL, MySQL, SQLite examples

## Database Services

The `docker-compose.yml` provides test databases:

- **PostgreSQL 16** - Port 5432
  - User: `postgres`
  - Password: `mysecretpassword`
  - Database: `postgres`

- **MySQL 9** - Port 3306
  - User: `root`
  - Password: `mysecretpassword`
  - Database: `dinotest`

## Makefile Targets

```bash
make test              # Basic unit tests only
make test-examples     # Tests with examples tag
make build-endtoend    # Build end-to-end test data
make test-ci           # Full CI suite (examples + endtoend + vet)
make vet               # Run go vet
make start             # Start database containers
```

## CI/CD Configuration

### GitHub Actions Workflow

- **File:** `.github/workflows/ci.yml`
- **Go Version:** 1.26.0
- **Database Setup:** Uses `sqlc-test-setup` (not Docker) to install and start PostgreSQL and MySQL directly on the runner
- **Test Command:** `gotestsum --junitfile junit.xml -- --tags=examples -timeout 20m ./...`
- **Additional Checks:** `govulncheck` for vulnerability scanning

## Development Workflow

### Building Development Versions

```bash
# Build main sqlc binary for development
go build -o ~/go/bin/sqlc-dev ./cmd/sqlc

# Build JSON plugin (required for some tests)
go build -o ~/go/bin/sqlc-gen-json ./cmd/sqlc-gen-json
```

### Environment Variables for Tests

You can override database connections via environment variables:

```bash
POSTGRESQL_SERVER_URI="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
MYSQL_SERVER_URI="root:mysecretpassword@tcp(127.0.0.1:3306)/mysql?multiStatements=true&parseTime=true"
```

## Code Structure

### Key Directories

- `/cmd/` - Main binaries (sqlc, sqlc-gen-json, sqlc-test-setup)
- `/internal/cmd/` - Command implementations (vet, generate, etc.)
- `/internal/engine/` - Database engine implementations
  - `/postgresql/` - PostgreSQL parser and converter
  - `/dolphin/` - MySQL parser (uses TiDB parser)
  - `/sqlite/` - SQLite parser
- `/internal/compiler/` - Query compilation logic
- `/internal/codegen/` - Code generation for different languages
- `/internal/config/` - Configuration file parsing
- `/internal/endtoend/` - End-to-end tests
- `/internal/sqltest/` - Test database setup (Docker, native, local detection)
- `/examples/` - Example projects for testing

### Important Files

- `/Makefile` - Build and test targets
- `/docker-compose.yml` - Database services for testing
- `/.github/workflows/ci.yml` - CI configuration

## Common Issues & Solutions

### Network Connectivity Issues

If you see errors about `storage.googleapis.com`, the Go proxy may be unreachable. Use `GOPROXY=direct go mod download` to fetch modules directly from source.

### Test Timeouts

End-to-end tests can take a while. Use longer timeouts:
```bash
go test -timeout 20m --tags=examples ./...
```

### Race Conditions

Always run tests with the race detector when working on concurrent code:
```bash
go test -race ./...
```

### Database Connection Failures

If using Docker:
```bash
docker compose ps
docker compose up -d
```

If using sqlc-test-setup:
```bash
go run ./cmd/sqlc-test-setup start
```

## Tips for Contributors

1. **Run tests before committing:** `go test --tags=examples -timeout 20m ./...`
2. **Check for race conditions:** Use `-race` flag when testing concurrent code
3. **Use specific package tests:** Faster iteration during development
4. **Read existing tests:** Good examples in `/internal/engine/postgresql/*_test.go`

## Git Workflow

### Branch Naming

- Feature branches should start with `claude/` for Claude Code work
- Branch names should be descriptive and end with the session ID

### Committing Changes

```bash
git add <files>
git commit -m "Brief description of changes"
git push -u origin <branch-name>
```

### Rebasing

```bash
git checkout main
git pull origin main
git checkout <feature-branch>
git rebase main
git push --force-with-lease origin <feature-branch>
```

## Resources

- **Main Documentation:** `/docs/`
- **Development Guide:** `/docs/guides/development.md`
- **CI Configuration:** `/.github/workflows/ci.yml`
- **Docker Compose:** `/docker-compose.yml`
