# Claude Code Development Guide for sqlc

This document provides essential information for working with the sqlc codebase, including testing, development workflow, and code structure.

## Quick Start

### Prerequisites

- **Go 1.25.0+** - Required for building and testing
- **Docker & Docker Compose** - Required for integration tests with databases
- **Git** - For version control

### Running Tests

#### Basic Unit Tests (No Database Required)

```bash
# Simplest approach - runs all unit tests
go test ./...

# Using make
make test
```

#### Full Test Suite with Integration Tests

```bash
# Step 1: Start database containers
docker compose up -d

# Step 2: Run all tests including examples
go test --tags=examples -timeout 20m ./...

# Or use make for the full CI suite
make test-ci
```

#### Running Specific Tests

```bash
# Test a specific package
go test ./internal/config
go test ./internal/compiler

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

### Managing Databases

```bash
# Start databases
make start
# or
docker compose up -d

# Stop databases
docker compose down

# View logs
docker compose logs -f
```

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
- **Go Version:** 1.25.0
- **Test Command:** `gotestsum --junitfile junit.xml -- --tags=examples -timeout 20m ./...`
- **Additional Checks:** `govulncheck` for vulnerability scanning

### Running Tests Like CI Locally

```bash
# Install CI tools (optional)
go install gotest.tools/gotestsum@latest

# Run tests with same timeout as CI
go test --tags=examples -timeout 20m ./...

# Or use the CI make target
make test-ci
```

## Development Workflow

### Building Development Versions

```bash
# Build main sqlc binary for development
go build -o ~/go/bin/sqlc-dev ./cmd/sqlc

# Build JSON plugin (required for some tests)
go build -o ~/go/bin/sqlc-gen-json ./cmd/sqlc-gen-json
```

### Environment Variables for Tests

You can customize database connections:

**PostgreSQL:**
```bash
PG_HOST=127.0.0.1
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=mysecretpassword
PG_DATABASE=dinotest
```

**MySQL:**
```bash
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_ROOT_PASSWORD=mysecretpassword
MYSQL_DATABASE=dinotest
```

**Example:**
```bash
POSTGRESQL_SERVER_URI="postgres://postgres:mysecretpassword@localhost:5432/postgres" \
  go test -v ./...
```

## Code Structure

### Key Directories

- `/cmd/` - Main binaries (sqlc, sqlc-gen-json)
- `/internal/cmd/` - Command implementations (vet, generate, etc.)
- `/internal/engine/` - Database engine implementations
  - `/clickhouse/` - ClickHouse parser and converter
  - `/postgresql/` - PostgreSQL parser and converter
  - `/dolphin/` - MySQL parser (uses TiDB parser)
  - `/sqlite/` - SQLite parser
- `/internal/compiler/` - Query compilation logic
- `/internal/codegen/` - Code generation for different languages
- `/internal/config/` - Configuration file parsing
- `/internal/endtoend/` - End-to-end tests
- `/examples/` - Example projects for testing

### Important Files

- `/Makefile` - Build and test targets
- `/docker-compose.yml` - Database services for testing
- `/.github/workflows/ci.yml` - CI configuration
- `/docs/guides/development.md` - Developer documentation

## Common Issues & Solutions

### Network Connectivity Issues

If you see errors about `storage.googleapis.com`, the Go proxy may be unreachable. Tests may still pass for packages that don't require network dependencies.

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

Ensure Docker containers are running:
```bash
docker compose ps
docker compose up -d
```

## Tips for Contributors

1. **Run tests before committing:** `make test-ci`
2. **Check for race conditions:** Use `-race` flag when testing concurrent code
3. **Use specific package tests:** Faster iteration during development
4. **Start databases early:** `docker compose up -d` before running integration tests
5. **Read existing tests:** Good examples in `/internal/engine/postgresql/*_test.go`

## Git Workflow

### Branch Naming

- Feature branches should start with `claude/` for Claude Code work
- Branch names should be descriptive and end with the session ID

### Committing Changes

```bash
# Stage changes
git add <files>

# Commit with descriptive message
git commit -m "Brief description

Detailed explanation of changes.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Push to remote
git push -u origin <branch-name>
```

### Rebasing

```bash
# Update main
git checkout main
git pull origin main

# Rebase feature branch
git checkout <feature-branch>
git rebase main

# Force push rebased branch
git push --force-with-lease origin <feature-branch>
```

## Resources

- **Main Documentation:** `/docs/`
- **Development Guide:** `/docs/guides/development.md`
- **CI Configuration:** `/.github/workflows/ci.yml`
- **Docker Compose:** `/docker-compose.yml`

## ClickHouse Engine Support

### Implementation Details

The ClickHouse engine was added to support ClickHouse SQL parsing and code generation using the `github.com/AfterShip/clickhouse-sql-parser` library.

**Files involved:**
- `internal/engine/clickhouse/parse.go` - Parser implementation
- `internal/engine/clickhouse/convert.go` - AST converter
- `internal/engine/clickhouse/catalog.go` - Catalog initialization
- `internal/config/config.go` - Engine constant registration
- `internal/compiler/engine.go` - Compiler integration
- `internal/compiler/expand.go` - Quote character handling
- `examples/clickhouse/` - Example project

**Key features:**
- Parses ClickHouse SQL into sqlc's internal AST
- Supports comment syntax: `--` and `/* */`
- Uses backticks for identifier quoting (ClickHouse standard)
- Handles ClickHouse-specific keywords

**Testing:**
```bash
# Test ClickHouse engine (unit tests)
go test ./internal/engine/clickhouse -v

# Test compiler integration
go test ./internal/compiler -v -run TestNewCompilerClickHouse

# Test code generation
cd examples/clickhouse && /path/to/sqlc generate && go mod init github.com/example/clickhouse && go mod tidy && go build ./gen
```

**Test Coverage:**
- Basic SELECT queries with WHERE, ORDER BY, LIMIT
- INSERT statements with VALUES
- Aggregate queries (COUNT, SUM, AVG) with GROUP BY and HAVING
- UNION and UNION ALL queries
- Subqueries and derived tables
- Multiple JOINs (INNER, LEFT, RIGHT, FULL)
- Window functions with OVER clause
- CAST expressions
- CASE expressions
- IS NULL / IS NOT NULL expressions
- Unary expressions (NOT, negation)
- Number literals (int, float, scientific notation)
- String literals and functions
- ClickHouse-specific features (PREWHERE, SAMPLE, array functions)

## Recent Fixes & Improvements

### Fixed Issues

1. **Typo in create_function_stmt.go** - Fixed "Undertand" → "Understand"
2. **Race condition in vet.go** - Fixed Client initialization using `sync.Once`
3. **Nil pointer dereference in parse.go** - Fixed unsafe type assertion in primary key parsing
4. **ClickHouse engine addition** - Added full ClickHouse SQL support with AST parsing

These fixes demonstrate common patterns:
- Using `sync.Once` for thread-safe lazy initialization
- Using comma-ok idiom for safe type assertions: `if val, ok := x.(Type); ok { ... }`
- Adding proper nil checks and defensive programming
- Registering new database engines following the established pattern

---

**Last Updated:** 2025-11-19
**Maintainer:** Claude Code
