# Claude Code Development Guide for sqlc

This document provides essential information for working with the sqlc codebase, including testing, development workflow, and code structure.

## Quick Start

### Prerequisites

- **Go 1.26.0+** - Required for building and testing
- **Docker & Docker Compose** - Required for integration tests with databases (local development)
- **Git** - For version control

## Claude Code Remote Environment Setup

When running in the Claude Code remote environment (or any environment without Docker), you can install PostgreSQL and MySQL natively. The test framework automatically detects and uses native database installations.

### Network Proxy (Pre-configured)

The Claude Code remote environment routes outbound traffic through an HTTP proxy via the `HTTP_PROXY` and `HTTPS_PROXY` environment variables. **Go module operations (`go mod tidy`, `go mod download`, `go get`, etc.) work automatically** because Go's toolchain respects these variables. No extra configuration is needed for the Go module proxy (`proxy.golang.org`) or checksum database (`sum.golang.org`).

### Step 1: Configure apt Proxy (Required in Remote Environment)

The apt package manager needs its own proxy configuration since it does not read `HTTP_PROXY`:

```bash
bash -c 'echo "Acquire::http::Proxy \"$http_proxy\";"' | sudo tee /etc/apt/apt.conf.d/99proxy
```

### Step 2: Install PostgreSQL

```bash
sudo apt-get update
sudo apt-get install -y postgresql
sudo service postgresql start
```

Configure PostgreSQL for password authentication:

```bash
# Set password for postgres user
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"

# Enable password authentication for localhost
echo 'host    all             all             127.0.0.1/32            md5' | sudo tee -a /etc/postgresql/16/main/pg_hba.conf
sudo service postgresql reload
```

Test the connection:

```bash
PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -c "SELECT 1;"
```

### Step 3: Install MySQL 9

MySQL 9 is required for full test compatibility (includes VECTOR type support). Download and install from Oracle:

```bash
# Download MySQL 9 bundle
curl -LO https://dev.mysql.com/get/Downloads/MySQL-9.1/mysql-server_9.1.0-1ubuntu24.04_amd64.deb-bundle.tar

# Extract packages
mkdir -p /tmp/mysql9
tar -xf mysql-server_9.1.0-1ubuntu24.04_amd64.deb-bundle.tar -C /tmp/mysql9

# Install packages (in order)
cd /tmp/mysql9
sudo dpkg -i mysql-common_*.deb \
    mysql-community-client-plugins_*.deb \
    mysql-community-client-core_*.deb \
    mysql-community-client_*.deb \
    mysql-client_*.deb \
    mysql-community-server-core_*.deb \
    mysql-community-server_*.deb \
    mysql-server_*.deb

# Make init script executable
sudo chmod +x /etc/init.d/mysql

# Initialize data directory and start MySQL
sudo mysqld --initialize-insecure --user=mysql
sudo /etc/init.d/mysql start

# Set root password
mysql -u root -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'mysecretpassword'; FLUSH PRIVILEGES;"
```

Test the connection:

```bash
mysql -h 127.0.0.1 -u root -pmysecretpassword -e "SELECT VERSION();"
```

### Step 4: Run End-to-End Tests

With both databases running, the test framework automatically detects them:

```bash
# Run all end-to-end tests
go test --tags=examples -timeout 20m ./internal/endtoend/...

# Run example tests
go test --tags=examples -timeout 20m ./examples/...

# Run the full test suite
go test --tags=examples -timeout 20m ./...
```

The native database support (in `internal/sqltest/native/`) automatically:
- Detects running PostgreSQL and MySQL instances
- Starts services if installed but not running
- Uses standard connection URIs:
  - PostgreSQL: `postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable`
  - MySQL: `root:mysecretpassword@tcp(127.0.0.1:3306)/mysql`

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

### Network Connectivity / Go Module Proxy Issues

In the Claude Code remote environment, Go module fetching works automatically via the `HTTP_PROXY`/`HTTPS_PROXY` environment variables. If you see errors about `storage.googleapis.com` or `proxy.golang.org`:

1. **Verify proxy vars are set:** `echo $HTTP_PROXY` (should be non-empty in the remote environment)
2. **Test connectivity:** `go mod download 2>&1 | head` to check for errors
3. **If proxy is not set** (e.g., local development without Docker), Go will try to reach `proxy.golang.org` directly, which requires internet access

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

## Recent Fixes & Improvements

### Fixed Issues

1. **Typo in create_function_stmt.go** - Fixed "Undertand" → "Understand"
2. **Race condition in vet.go** - Fixed Client initialization using `sync.Once`
3. **Nil pointer dereference in parse.go** - Fixed unsafe type assertion in primary key parsing

These fixes demonstrate common patterns:
- Using `sync.Once` for thread-safe lazy initialization
- Using comma-ok idiom for safe type assertions: `if val, ok := x.(Type); ok { ... }`
- Adding proper nil checks and defensive programming

---

**Last Updated:** 2025-10-21
**Maintainer:** Claude Code
