package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	log.SetFlags(log.Ltime)
	log.SetPrefix("[sqlc-test-setup] ")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: sqlc-test-setup <install|start>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		if err := runInstall(); err != nil {
			log.Fatalf("install failed: %s", err)
		}
	case "start":
		if err := runStart(); err != nil {
			log.Fatalf("start failed: %s", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\nusage: sqlc-test-setup <install|start>\n", os.Args[1])
		os.Exit(1)
	}
}

// run executes a command with verbose logging, streaming output to stderr.
func run(name string, args ...string) error {
	log.Printf("exec: %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// runOutput executes a command and returns its combined output.
func runOutput(name string, args ...string) (string, error) {
	log.Printf("exec: %s %s", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runInstall() error {
	log.Println("=== Installing PostgreSQL and MySQL for test setup ===")

	if err := installAptProxy(); err != nil {
		return fmt.Errorf("configuring apt proxy: %w", err)
	}

	if err := installPostgreSQL(); err != nil {
		return fmt.Errorf("installing postgresql: %w", err)
	}

	if err := installMySQL(); err != nil {
		return fmt.Errorf("installing mysql: %w", err)
	}

	log.Println("=== Install complete ===")
	return nil
}

func installAptProxy() error {
	proxy := os.Getenv("http_proxy")
	if proxy == "" {
		log.Println("http_proxy is not set, skipping apt proxy configuration")
		return nil
	}

	log.Printf("configuring apt proxy to use %s", proxy)
	proxyConf := fmt.Sprintf("Acquire::http::Proxy \"%s\";", proxy)
	cmd := fmt.Sprintf("echo '%s' | sudo tee /etc/apt/apt.conf.d/99proxy", proxyConf)
	return run("bash", "-c", cmd)
}

func installPostgreSQL() error {
	log.Println("--- Installing PostgreSQL ---")

	log.Println("updating apt package lists")
	if err := run("sudo", "apt-get", "update", "-qq"); err != nil {
		return fmt.Errorf("apt-get update: %w", err)
	}

	log.Println("installing postgresql package")
	if err := run("sudo", "apt-get", "install", "-y", "-qq", "postgresql"); err != nil {
		return fmt.Errorf("apt-get install postgresql: %w", err)
	}

	log.Println("postgresql installed successfully")
	return nil
}

func installMySQL() error {
	log.Println("--- Installing MySQL 9 ---")

	bundleURL := "https://dev.mysql.com/get/Downloads/MySQL-9.1/mysql-server_9.1.0-1ubuntu24.04_amd64.deb-bundle.tar"
	bundleTar := "/tmp/mysql-server-bundle.tar"
	extractDir := "/tmp/mysql9"

	log.Printf("downloading MySQL 9 bundle from %s", bundleURL)
	if err := run("curl", "-L", "-o", bundleTar, bundleURL); err != nil {
		return fmt.Errorf("downloading mysql bundle: %w", err)
	}

	log.Printf("extracting bundle to %s", extractDir)
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("creating extract dir: %w", err)
	}
	if err := run("tar", "-xf", bundleTar, "-C", extractDir); err != nil {
		return fmt.Errorf("extracting mysql bundle: %w", err)
	}

	// Install packages in dependency order
	packages := []string{
		"mysql-common_*.deb",
		"mysql-community-client-plugins_*.deb",
		"mysql-community-client-core_*.deb",
		"mysql-community-client_*.deb",
		"mysql-client_*.deb",
		"mysql-community-server-core_*.deb",
		"mysql-community-server_*.deb",
		"mysql-server_*.deb",
	}

	for _, pkg := range packages {
		log.Printf("installing %s", pkg)
		// Use shell glob expansion via bash -c
		cmd := fmt.Sprintf("sudo dpkg -i %s/%s", extractDir, pkg)
		if err := run("bash", "-c", cmd); err != nil {
			return fmt.Errorf("installing %s: %w", pkg, err)
		}
	}

	log.Println("making mysql init script executable")
	if err := run("sudo", "chmod", "+x", "/etc/init.d/mysql"); err != nil {
		return fmt.Errorf("chmod mysql init script: %w", err)
	}

	log.Println("mysql 9 installed successfully")
	return nil
}

func runStart() error {
	log.Println("=== Starting PostgreSQL and MySQL ===")

	if err := startPostgreSQL(); err != nil {
		return fmt.Errorf("starting postgresql: %w", err)
	}

	if err := startMySQL(); err != nil {
		return fmt.Errorf("starting mysql: %w", err)
	}

	log.Println("=== Both databases are running and configured ===")
	log.Println("PostgreSQL: postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	log.Println("MySQL:      root:mysecretpassword@tcp(127.0.0.1:3306)/mysql")
	return nil
}

func startPostgreSQL() error {
	log.Println("--- Starting PostgreSQL ---")

	log.Println("starting postgresql service")
	if err := run("sudo", "service", "postgresql", "start"); err != nil {
		return fmt.Errorf("service postgresql start: %w", err)
	}

	log.Println("setting password for postgres user")
	if err := run("sudo", "-u", "postgres", "psql", "-c", "ALTER USER postgres PASSWORD 'postgres';"); err != nil {
		return fmt.Errorf("setting postgres password: %w", err)
	}

	log.Println("detecting postgresql config directory")
	hbaPath, err := detectPgHBAPath()
	if err != nil {
		return fmt.Errorf("detecting pg_hba.conf path: %w", err)
	}

	log.Printf("enabling md5 authentication in %s", hbaPath)
	hbaLine := "host    all             all             127.0.0.1/32            md5"
	cmd := fmt.Sprintf("echo '%s' | sudo tee -a %s", hbaLine, hbaPath)
	if err := run("bash", "-c", cmd); err != nil {
		return fmt.Errorf("configuring pg_hba.conf: %w", err)
	}

	log.Println("reloading postgresql configuration")
	if err := run("sudo", "service", "postgresql", "reload"); err != nil {
		return fmt.Errorf("reloading postgresql: %w", err)
	}

	log.Println("verifying postgresql connection")
	if err := run("bash", "-c", "PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -c 'SELECT 1;'"); err != nil {
		return fmt.Errorf("postgresql connection test failed: %w", err)
	}

	log.Println("postgresql is running and configured")
	return nil
}

// detectPgHBAPath finds the pg_hba.conf file across different PostgreSQL versions.
func detectPgHBAPath() (string, error) {
	out, err := runOutput("bash", "-c", "sudo -u postgres psql -t -c 'SHOW hba_file;'")
	if err != nil {
		return "", fmt.Errorf("querying hba_file: %w (output: %s)", err, out)
	}
	path := strings.TrimSpace(out)
	if path == "" {
		return "", fmt.Errorf("pg_hba.conf path is empty")
	}
	log.Printf("found pg_hba.conf at %s", path)
	return path, nil
}

func startMySQL() error {
	log.Println("--- Starting MySQL ---")

	log.Println("initializing mysql data directory")
	if err := run("sudo", "mysqld", "--initialize-insecure", "--user=mysql"); err != nil {
		return fmt.Errorf("mysqld --initialize-insecure: %w", err)
	}

	log.Println("starting mysql service")
	if err := run("sudo", "/etc/init.d/mysql", "start"); err != nil {
		return fmt.Errorf("starting mysql: %w", err)
	}

	log.Println("setting mysql root password")
	if err := run("mysql", "-u", "root", "-e",
		"ALTER USER 'root'@'localhost' IDENTIFIED BY 'mysecretpassword'; FLUSH PRIVILEGES;"); err != nil {
		return fmt.Errorf("setting mysql root password: %w", err)
	}

	log.Println("verifying mysql connection")
	if err := run("mysql", "-h", "127.0.0.1", "-u", "root", "-pmysecretpassword", "-e", "SELECT VERSION();"); err != nil {
		return fmt.Errorf("mysql connection test failed: %w", err)
	}

	log.Println("mysql is running and configured")
	return nil
}
