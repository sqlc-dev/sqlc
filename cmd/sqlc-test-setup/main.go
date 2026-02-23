package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// pgVersion is the PostgreSQL version to install.
	pgVersion = "18.2.0"
)

// pgBinary contains the download information for a PostgreSQL binary release.
type pgBinary struct {
	URL    string
	SHA256 string
}

// pgBinaries maps "<GOOS>/<GOARCH>" to the corresponding binary download info.
var pgBinaries = map[string]pgBinary{
	"linux/amd64": {
		URL:    "https://github.com/theseus-rs/postgresql-binaries/releases/download/" + pgVersion + "/postgresql-" + pgVersion + "-x86_64-unknown-linux-gnu.tar.gz",
		SHA256: "cc2674e1641aa2a62b478971a22c131a768eb783f313e6a3385888f58a604074",
	},
	"linux/arm64": {
		URL:    "https://github.com/theseus-rs/postgresql-binaries/releases/download/" + pgVersion + "/postgresql-" + pgVersion + "-aarch64-unknown-linux-gnu.tar.gz",
		SHA256: "8b415a11c7a5484e5fbf7a57fca71554d2d1d7acd34faf066606d2fee1261854",
	},
}

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

// pgBaseDir returns the sqlc-specific directory where PostgreSQL is installed,
// using the user's cache directory (~/.cache/sqlc/postgresql on Linux).
func pgBaseDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return filepath.Join(cacheDir, "sqlc", "postgresql")
}

// pgBinDir returns the path to the PostgreSQL bin directory.
func pgBinDir() string {
	return filepath.Join(pgBaseDir(), "bin")
}

// pgDataDir returns the path to the PostgreSQL data directory.
func pgDataDir() string {
	return filepath.Join(pgBaseDir(), "data")
}

// pgBin returns the full path to a PostgreSQL binary.
func pgBin(name string) string {
	return filepath.Join(pgBinDir(), name)
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

	// Install runtime dependencies needed by PostgreSQL extensions (e.g.
	// uuid-ossp requires libossp-uuid16).
	if err := installPgDeps(); err != nil {
		return fmt.Errorf("installing postgresql dependencies: %w", err)
	}

	// Check if already installed in our directory
	if _, err := os.Stat(pgBin("postgres")); err == nil {
		out, err := runOutput(pgBin("postgres"), "--version")
		if err == nil {
			log.Printf("postgresql is already installed: %s", strings.TrimSpace(out))
			log.Println("skipping postgresql installation")
			return nil
		}
	}

	platform := runtime.GOOS + "/" + runtime.GOARCH
	bin, ok := pgBinaries[platform]
	if !ok {
		return fmt.Errorf("unsupported platform: %s (supported: %s)", platform, supportedPlatforms())
	}

	// Download to a temp file
	tarball := filepath.Join(os.TempDir(), fmt.Sprintf("postgresql-%s.tar.gz", pgVersion))

	if _, err := os.Stat(tarball); err != nil {
		log.Printf("downloading PostgreSQL %s from %s", pgVersion, bin.URL)
		if err := downloadFile(tarball, bin.URL); err != nil {
			os.Remove(tarball)
			return fmt.Errorf("downloading postgresql: %w", err)
		}
	} else {
		log.Printf("postgresql tarball already downloaded at %s", tarball)
	}

	// Verify SHA256 checksum
	log.Printf("verifying SHA256 checksum")
	actualHash, err := sha256File(tarball)
	if err != nil {
		return fmt.Errorf("computing sha256: %w", err)
	}
	if actualHash != bin.SHA256 {
		os.Remove(tarball)
		return fmt.Errorf("SHA256 mismatch: expected %s, got %s", bin.SHA256, actualHash)
	}
	log.Printf("SHA256 checksum verified: %s", actualHash)

	baseDir := pgBaseDir()

	// Create the base directory in the user cache
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", baseDir, err)
	}

	// Extract the tarball - it contains a top-level directory like
	// postgresql-18.2.0-x86_64-unknown-linux-gnu/ with bin/, lib/, share/ inside.
	// We strip that top-level directory and extract directly into the base dir.
	log.Printf("extracting postgresql to %s", baseDir)
	if err := run("tar", "-xzf", tarball, "-C", baseDir, "--strip-components=1"); err != nil {
		return fmt.Errorf("extracting postgresql: %w", err)
	}

	// Verify the binary works
	out, err := runOutput(pgBin("postgres"), "--version")
	if err != nil {
		return fmt.Errorf("postgres --version failed after install: %w", err)
	}
	log.Printf("postgresql installed successfully: %s", strings.TrimSpace(out))
	return nil
}

// installPgDeps installs shared libraries required by PostgreSQL extensions at
// runtime (e.g. libossp-uuid16 for uuid-ossp).
func installPgDeps() error {
	log.Println("installing postgresql runtime dependencies")
	if err := run("sudo", "apt-get", "install", "-y", "--no-install-recommends", "libossp-uuid16"); err != nil {
		return fmt.Errorf("apt-get install libossp-uuid16: %w", err)
	}
	return nil
}

// supportedPlatforms returns a comma-separated list of supported platforms.
func supportedPlatforms() string {
	platforms := make([]string, 0, len(pgBinaries))
	for p := range pgBinaries {
		platforms = append(platforms, p)
	}
	return strings.Join(platforms, ", ")
}

// downloadFile downloads a URL to a local file path.
func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// sha256File computes the SHA256 hash of a file and returns the hex string.
func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
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

	dataDir := pgDataDir()
	logFile := filepath.Join(pgBaseDir(), "postgresql.log")

	// Check if already running
	if pgIsReady() {
		log.Println("postgresql is already running and accepting connections")
		return nil
	}

	// Initialize data directory if needed
	if _, err := os.Stat(filepath.Join(dataDir, "PG_VERSION")); os.IsNotExist(err) {
		log.Println("initializing postgresql data directory")
		if err := os.MkdirAll(dataDir, 0o700); err != nil {
			return fmt.Errorf("creating data directory: %w", err)
		}
		if err := run(pgBin("initdb"),
			"-D", dataDir,
			"--username=postgres",
			"--auth=trust",
		); err != nil {
			return fmt.Errorf("initdb: %w", err)
		}

		// Configure pg_hba.conf for md5 password authentication on TCP
		hbaPath := filepath.Join(dataDir, "pg_hba.conf")
		if err := configurePgHBA(hbaPath); err != nil {
			return fmt.Errorf("configuring pg_hba.conf: %w", err)
		}

		// Configure postgresql.conf to listen on localhost
		confPath := filepath.Join(dataDir, "postgresql.conf")
		if err := appendToFile(confPath,
			"\n# sqlc-test-setup configuration\n"+
				"listen_addresses = '127.0.0.1'\n"+
				"port = 5432\n",
		); err != nil {
			return fmt.Errorf("configuring postgresql.conf: %w", err)
		}
	} else {
		log.Println("postgresql data directory already initialized")
	}

	// Start PostgreSQL using pg_ctl
	log.Println("starting postgresql")
	if err := run(pgBin("pg_ctl"),
		"-D", dataDir,
		"-l", logFile,
		"-o", fmt.Sprintf("-k %s", dataDir),
		"start",
	); err != nil {
		return fmt.Errorf("pg_ctl start: %w", err)
	}

	// Wait for PostgreSQL to be ready
	log.Println("waiting for postgresql to accept connections")
	if err := waitForPostgreSQL(30 * time.Second); err != nil {
		return fmt.Errorf("postgresql did not start in time: %w", err)
	}

	// Set the postgres user password
	log.Println("setting password for postgres user")
	if err := run(pgBin("psql"),
		"-h", "127.0.0.1",
		"-U", "postgres",
		"-c", "ALTER USER postgres PASSWORD 'postgres';",
	); err != nil {
		return fmt.Errorf("setting postgres password: %w", err)
	}

	// Update pg_hba.conf to require md5 auth now that password is set
	hbaPath := filepath.Join(dataDir, "pg_hba.conf")
	if err := configurePgHBAWithMD5(hbaPath); err != nil {
		return fmt.Errorf("updating pg_hba.conf for md5: %w", err)
	}

	// Reload configuration
	log.Println("reloading postgresql configuration")
	if err := run(pgBin("pg_ctl"), "-D", dataDir, "reload"); err != nil {
		return fmt.Errorf("pg_ctl reload: %w", err)
	}

	// Verify connection with password
	log.Println("verifying postgresql connection")
	cmd := exec.Command(pgBin("psql"),
		"-h", "127.0.0.1",
		"-U", "postgres",
		"-c", "SELECT 1;",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD=postgres")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("postgresql connection test failed: %w", err)
	}

	log.Println("postgresql is running and configured")
	return nil
}

// configurePgHBA writes a pg_hba.conf that allows trust auth initially (for
// setting the password), then we switch to md5.
func configurePgHBA(hbaPath string) error {
	content := `# pg_hba.conf - generated by sqlc-test-setup
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            trust
host    all             all             ::1/128                 trust
`
	return os.WriteFile(hbaPath, []byte(content), 0o600)
}

// configurePgHBAWithMD5 rewrites pg_hba.conf to use md5 for TCP connections.
func configurePgHBAWithMD5(hbaPath string) error {
	content := `# pg_hba.conf - generated by sqlc-test-setup
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5
`
	return os.WriteFile(hbaPath, []byte(content), 0o600)
}

// appendToFile appends text to a file.
func appendToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(text)
	return err
}

// pgIsReady checks if PostgreSQL is running and accepting connections.
func pgIsReady() bool {
	cmd := exec.Command(pgBin("pg_isready"), "-h", "127.0.0.1", "-p", "5432")
	return cmd.Run() == nil
}

// waitForPostgreSQL polls until PostgreSQL accepts connections or times out.
func waitForPostgreSQL(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if pgIsReady() {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out after %s waiting for postgresql", timeout)
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
