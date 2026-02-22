package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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

// commandExists checks if a binary is available in PATH.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// isMySQLVersionOK checks if the mysqld --version output indicates MySQL 9+.
// Example version string: "/usr/sbin/mysqld  Ver 8.0.44-0ubuntu0.24.04.2 ..."
func isMySQLVersionOK(versionOutput string) bool {
	// Look for "Ver X.Y.Z" pattern
	fields := strings.Fields(versionOutput)
	for i, f := range fields {
		if strings.EqualFold(f, "Ver") && i+1 < len(fields) {
			ver := strings.Split(fields[i+1], ".")
			if len(ver) > 0 {
				major := strings.TrimLeft(ver[0], "0")
				if major == "" {
					return false
				}
				return major[0] >= '9'
			}
		}
	}
	return false
}

// ---- install ----

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

	const confPath = "/etc/apt/apt.conf.d/99proxy"
	if _, err := os.Stat(confPath); err == nil {
		log.Printf("apt proxy config already exists at %s, skipping", confPath)
		return nil
	}

	log.Printf("configuring apt proxy to use %s", proxy)
	proxyConf := fmt.Sprintf("Acquire::http::Proxy \"%s\";", proxy)
	cmd := fmt.Sprintf("echo '%s' | sudo tee /etc/apt/apt.conf.d/99proxy", proxyConf)
	return run("bash", "-c", cmd)
}

func installPostgreSQL() error {
	log.Println("--- Installing PostgreSQL ---")

	if commandExists("psql") {
		out, err := runOutput("psql", "--version")
		if err == nil {
			log.Printf("postgresql is already installed: %s", strings.TrimSpace(out))
			log.Println("skipping postgresql installation")
			return nil
		}
	}

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

	if commandExists("mysqld") {
		out, err := runOutput("mysqld", "--version")
		if err == nil {
			version := strings.TrimSpace(out)
			log.Printf("mysql is already installed: %s", version)
			if isMySQLVersionOK(version) {
				log.Println("mysql version is 9+, skipping installation")
				return nil
			}
			log.Println("mysql version is too old, upgrading to MySQL 9")
			// Stop existing MySQL before upgrading
			_ = exec.Command("sudo", "service", "mysql", "stop").Run()
			_ = exec.Command("sudo", "pkill", "-f", "mysqld").Run()
			time.Sleep(2 * time.Second)
			// Remove old MySQL packages to avoid conflicts
			log.Println("removing old mysql packages")
			_ = run("sudo", "apt-get", "remove", "-y", "mysql-server", "mysql-client", "mysql-common",
				"mysql-server-core-*", "mysql-client-core-*")
			// Clear old data directory so MySQL 9 can initialize fresh
			log.Println("clearing old mysql data directory")
			_ = run("sudo", "rm", "-rf", "/var/lib/mysql")
			_ = run("sudo", "mkdir", "-p", "/var/lib/mysql")
			_ = run("sudo", "chown", "mysql:mysql", "/var/lib/mysql")
		}
	}

	bundleURL := "https://dev.mysql.com/get/Downloads/MySQL-9.1/mysql-server_9.1.0-1ubuntu24.04_amd64.deb-bundle.tar"
	bundleTar := "/tmp/mysql-server-bundle.tar"
	extractDir := "/tmp/mysql9"

	if _, err := os.Stat(bundleTar); err != nil {
		log.Printf("downloading MySQL 9 bundle from %s", bundleURL)
		if err := run("curl", "-L", "-o", bundleTar, bundleURL); err != nil {
			return fmt.Errorf("downloading mysql bundle: %w", err)
		}
	} else {
		log.Printf("mysql bundle already downloaded at %s, skipping download", bundleTar)
	}

	log.Printf("extracting bundle to %s", extractDir)
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("creating extract dir: %w", err)
	}
	if err := run("tar", "-xf", bundleTar, "-C", extractDir); err != nil {
		return fmt.Errorf("extracting mysql bundle: %w", err)
	}

	// Install packages in dependency order using dpkg.
	// Some packages may fail due to missing dependencies, which is expected.
	// We fix them all at the end with apt-get install -f.
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
		log.Printf("installing %s (dependency errors will be fixed afterwards)", pkg)
		cmd := fmt.Sprintf("sudo dpkg -i %s/%s", extractDir, pkg)
		if err := run("bash", "-c", cmd); err != nil {
			log.Printf("dpkg reported errors for %s (will fix with apt-get install -f)", pkg)
		}
	}

	log.Println("fixing missing dependencies with apt-get install -f")
	if err := run("sudo", "apt-get", "install", "-f", "-y"); err != nil {
		return fmt.Errorf("apt-get install -f: %w", err)
	}

	log.Println("mysql 9 installed successfully")
	return nil
}

// ---- start ----

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

	if err := ensurePgHBAEntry(hbaPath); err != nil {
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

// ensurePgHBAEntry adds the md5 auth line to pg_hba.conf if it's not already present.
func ensurePgHBAEntry(hbaPath string) error {
	hbaLine := "host    all             all             127.0.0.1/32            md5"

	out, err := runOutput("sudo", "cat", hbaPath)
	if err != nil {
		return fmt.Errorf("reading pg_hba.conf: %w", err)
	}

	if strings.Contains(out, "127.0.0.1/32            md5") {
		log.Println("md5 authentication for 127.0.0.1/32 already configured in pg_hba.conf, skipping")
		return nil
	}

	log.Printf("enabling md5 authentication in %s", hbaPath)
	cmd := fmt.Sprintf("echo '%s' | sudo tee -a %s", hbaLine, hbaPath)
	return run("bash", "-c", cmd)
}

func startMySQL() error {
	log.Println("--- Starting MySQL ---")

	// Check if MySQL is already running and accessible with the expected password
	if mysqlReady() {
		log.Println("mysql is already running and accepting connections")
		return verifyMySQL()
	}

	// Stop any existing MySQL service that might be running (e.g. pre-installed
	// on GitHub Actions runners) to avoid port conflicts.
	log.Println("stopping any existing mysql service")
	_ = exec.Command("sudo", "service", "mysql", "stop").Run()
	_ = exec.Command("sudo", "mysqladmin", "shutdown").Run()
	// Give MySQL time to fully shut down
	time.Sleep(2 * time.Second)

	if err := ensureMySQLDirs(); err != nil {
		return err
	}

	// Check if data directory already exists and has been initialized
	needsPasswordReset := false
	if mysqlInitialized() {
		log.Println("mysql data directory already initialized, skipping initialization")
		// Existing data dir may have an unknown root password (e.g. pre-installed
		// MySQL on GitHub Actions). We'll need to use --skip-grant-tables to reset it.
		needsPasswordReset = true
	} else {
		log.Println("initializing mysql data directory")
		if err := run("sudo", "mysqld", "--initialize-insecure", "--user=mysql"); err != nil {
			return fmt.Errorf("mysqld --initialize-insecure: %w", err)
		}
	}

	if needsPasswordReset {
		// Start with --skip-grant-tables to reset the unknown root password.
		if err := startMySQLDaemon("--skip-grant-tables"); err != nil {
			return err
		}

		log.Println("resetting root password via --skip-grant-tables")
		resetSQL := "FLUSH PRIVILEGES; ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'mysecretpassword';"
		if err := run("mysql", "-u", "root", "-e", resetSQL); err != nil {
			return fmt.Errorf("resetting mysql root password: %w", err)
		}

		// Restart without --skip-grant-tables
		log.Println("restarting mysql normally")
		if err := run("sudo", "mysqladmin", "-u", "root", "-pmysecretpassword", "shutdown"); err != nil {
			// If mysqladmin fails, try killing the process directly
			_ = run("sudo", "pkill", "-f", "mysqld")
		}
		time.Sleep(2 * time.Second)

		if err := startMySQLDaemon(); err != nil {
			return err
		}
	} else {
		// Fresh initialization — start normally and set password
		if err := startMySQLDaemon(); err != nil {
			return err
		}

		log.Println("setting mysql root password")
		alterSQL := "ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'mysecretpassword'; FLUSH PRIVILEGES;"
		if err := run("mysql", "-u", "root", "-e", alterSQL); err != nil {
			return fmt.Errorf("setting mysql root password: %w", err)
		}
	}

	return verifyMySQL()
}

// ensureMySQLDirs creates the directories MySQL needs at runtime.
func ensureMySQLDirs() error {
	if err := run("sudo", "mkdir", "-p", "/var/run/mysqld"); err != nil {
		return fmt.Errorf("creating /var/run/mysqld: %w", err)
	}
	if err := run("sudo", "chown", "mysql:mysql", "/var/run/mysqld"); err != nil {
		return fmt.Errorf("chowning /var/run/mysqld: %w", err)
	}
	return nil
}

// startMySQLDaemon starts mysqld_safe in the background and waits for it to
// accept connections. Extra args (e.g. "--skip-grant-tables") are appended.
func startMySQLDaemon(extraArgs ...string) error {
	args := append([]string{"mysqld_safe", "--user=mysql"}, extraArgs...)
	log.Printf("starting mysql via mysqld_safe %v", extraArgs)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting mysqld_safe: %w", err)
	}

	log.Println("waiting for mysql to accept connections")
	if err := waitForMySQL(30 * time.Second); err != nil {
		return fmt.Errorf("mysql did not start in time: %w", err)
	}
	log.Println("mysql is accepting connections")
	return nil
}

// mysqlReady checks if MySQL is running and accepting connections with the expected password.
func mysqlReady() bool {
	err := exec.Command("mysqladmin", "-h", "127.0.0.1", "-u", "root", "-pmysecretpassword", "ping").Run()
	return err == nil
}

// waitForMySQL polls until MySQL accepts connections or the timeout expires.
func waitForMySQL(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// Try connecting without password (fresh) or with password (already configured)
		if exec.Command("mysqladmin", "-u", "root", "ping").Run() == nil {
			return nil
		}
		if exec.Command("mysqladmin", "-h", "127.0.0.1", "-u", "root", "-pmysecretpassword", "ping").Run() == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out after %s waiting for mysql", timeout)
}

func verifyMySQL() error {
	log.Println("verifying mysql connection")
	if err := run("mysql", "-h", "127.0.0.1", "-u", "root", "-pmysecretpassword", "-e", "SELECT VERSION();"); err != nil {
		return fmt.Errorf("mysql connection test failed: %w", err)
	}
	log.Println("mysql is running and configured")
	return nil
}

// mysqlInitialized checks if the MySQL data directory has been initialized.
// We use sudo ls because /var/lib/mysql is typically only readable by the
// mysql user, so filepath.Glob from a non-root process would silently fail.
func mysqlInitialized() bool {
	out, err := exec.Command("sudo", "ls", "/var/lib/mysql").CombinedOutput()
	if err != nil {
		return false
	}
	// If the directory has any contents, consider it initialized.
	// mysqld --initialize-insecure requires an empty directory.
	return strings.TrimSpace(string(out)) != ""
}
