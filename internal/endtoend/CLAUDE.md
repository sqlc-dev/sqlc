# End-to-End Tests - Native Database Setup

This document describes how to set up MySQL and PostgreSQL for running end-to-end tests in environments without Docker, particularly when using an HTTP proxy.

## Overview

The end-to-end tests support three methods for connecting to databases:

1. **Environment Variables**: Set `POSTGRESQL_SERVER_URI` and `MYSQL_SERVER_URI` directly
2. **Docker**: Automatically starts containers via the docker package
3. **Native Installation**: Starts existing database services on Linux

## Installing Databases with HTTP Proxy

In environments where DNS doesn't work directly but an HTTP proxy is available (e.g., some CI environments), you need to configure apt to use the proxy before installing packages.

### Configure apt Proxy

```bash
# Check if HTTP_PROXY is set
echo $HTTP_PROXY

# Configure apt to use the proxy
sudo tee /etc/apt/apt.conf.d/99proxy << EOF
Acquire::http::Proxy "$HTTP_PROXY";
Acquire::https::Proxy "$HTTPS_PROXY";
EOF

# Update package lists
sudo apt-get update -qq
```

### Install PostgreSQL

```bash
# Install PostgreSQL
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y postgresql postgresql-contrib

# Start the service
sudo service postgresql start

# Set password for postgres user
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"

# Configure pg_hba.conf for password authentication
# Find the hba_file location:
sudo -u postgres psql -t -c "SHOW hba_file;"

# Add md5 authentication for localhost (add to the beginning of pg_hba.conf):
# host    all             all             127.0.0.1/32            md5

# Reload PostgreSQL
sudo service postgresql reload
```

### Install MySQL

```bash
# Pre-configure MySQL root password
echo "mysql-server mysql-server/root_password password mysecretpassword" | sudo debconf-set-selections
echo "mysql-server mysql-server/root_password_again password mysecretpassword" | sudo debconf-set-selections

# Install MySQL
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y mysql-server

# Start the service
sudo service mysql start

# Verify connection
mysql -uroot -pmysecretpassword -e "SELECT 1;"
```

## Expected Database Credentials

The native database support expects the following credentials:

### PostgreSQL
- **URI**: `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`
- **User**: `postgres`
- **Password**: `postgres`
- **Port**: `5432`

### MySQL
- **URI**: `root:mysecretpassword@tcp(localhost:3306)/mysql?multiStatements=true&parseTime=true`
- **User**: `root`
- **Password**: `mysecretpassword`
- **Port**: `3306`

## Running Tests

```bash
# Run end-to-end tests
go test -v -run TestReplay -timeout 20m ./internal/endtoend/...

# With verbose logging
go test -v -run TestReplay -timeout 20m ./internal/endtoend/... 2>&1 | tee test.log
```

## Troubleshooting

### apt-get times out or fails
- Ensure HTTP proxy is configured in `/etc/apt/apt.conf.d/99proxy`
- Check that the proxy URL is correct: `echo $HTTP_PROXY`
- Try running `sudo apt-get update` first to verify connectivity

### MySQL connection refused
- Check if MySQL is running: `sudo service mysql status`
- Verify the password: `mysql -uroot -pmysecretpassword -e "SELECT 1;"`
- Check if MySQL is listening on TCP: `netstat -tlnp | grep 3306`

### PostgreSQL authentication failed
- Verify pg_hba.conf has md5 authentication for localhost
- Check password: `PGPASSWORD=postgres psql -h localhost -U postgres -c "SELECT 1;"`
- Reload PostgreSQL after config changes: `sudo service postgresql reload`

### DNS resolution fails
This is expected in some environments. Configure apt proxy as shown above.
