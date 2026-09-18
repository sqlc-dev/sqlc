package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// pgVersion is the PostgreSQL version to install.
	pgVersion = "18.2.0"

	// omniVersion is the Spanner Omni release to install, from its
	// standalone server binaries.
	omniVersion = "2026.r2.1-beta"

	// omniPort is the gRPC port a Spanner Omni single server listens on.
	omniPort = "15000"
)

// omniBinary contains the download information for a Spanner Omni server
// release, published at https://storage.googleapis.com/spanner-omni/ and
// documented at https://docs.cloud.google.com/spanner-omni/download.
type omniBinary struct {
	URL    string
	SHA256 string
}

// omniBinaries maps "<GOOS>/<GOARCH>" to the server download. Only linux
// x86_64 is published; the other platforms skip Spanner Omni.
var omniBinaries = map[string]omniBinary{
	"linux/amd64": {
		URL:    "https://storage.googleapis.com/spanner-omni/" + omniVersion + "/spanner-omni-server-" + omniVersion + "-linux-x86_64.tar.gz",
		SHA256: "792ffc772d5fff8ade56a8f336993d671d22794fe9f7e88e38a6609c65e16d90",
	},
}

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
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	services, err := selectServices(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n%s\n", err, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		if err := runInstall(services); err != nil {
			log.Fatalf("install failed: %s", err)
		}
	case "start":
		if err := runStart(services); err != nil {
			log.Fatalf("start failed: %s", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n%s\n", os.Args[1], usage)
		os.Exit(1)
	}
}

const usage = "usage: sqlc-test-setup <install|start> [postgresql|mysql|spanner ...]"

// allServices is every database the tool knows, in the order they are
// installed and started.
var allServices = []string{"postgresql", "mysql", "spanner"}

// selectServices reads the databases named on the command line, or every
// database when none is named.
func selectServices(args []string) (map[string]bool, error) {
	selected := map[string]bool{}
	if len(args) == 0 {
		for _, name := range allServices {
			selected[name] = true
		}
		return selected, nil
	}
	for _, arg := range args {
		known := false
		for _, name := range allServices {
			if arg == name {
				known = true
			}
		}
		if !known {
			return nil, fmt.Errorf("unknown database: %s", arg)
		}
		selected[arg] = true
	}
	return selected, nil
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

// The MySQL release the tests run against. Everything that names a MySQL
// version — the docker-compose service, the goldeneye dialect generator and
// its gen workflow — is pinned to the same release.
const (
	mysqlMajor   = 26
	mysqlVersion = "26.7.0"
)

// isMySQLVersionOK checks if the mysqld --version output indicates the
// pinned major release or a later one.
// Example version string: "/usr/sbin/mysqld  Ver 8.0.44-0ubuntu0.24.04.2 ..."
func isMySQLVersionOK(versionOutput string) bool {
	// Look for "Ver X.Y.Z" pattern
	fields := strings.Fields(versionOutput)
	for i, f := range fields {
		if strings.EqualFold(f, "Ver") && i+1 < len(fields) {
			ver := strings.Split(fields[i+1], ".")
			if len(ver) > 0 {
				major, err := strconv.Atoi(ver[0])
				if err != nil {
					return false
				}
				return major >= mysqlMajor
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

func runInstall(services map[string]bool) error {
	log.Println("=== Installing databases for test setup ===")

	if services["postgresql"] || services["mysql"] {
		if err := installAptProxy(); err != nil {
			return fmt.Errorf("configuring apt proxy: %w", err)
		}
	}

	if services["postgresql"] {
		if err := installPostgreSQL(); err != nil {
			return fmt.Errorf("installing postgresql: %w", err)
		}
	}

	if services["mysql"] {
		if err := installMySQL(); err != nil {
			return fmt.Errorf("installing mysql: %w", err)
		}
	}

	if services["spanner"] {
		if err := installSpannerOmni(); err != nil {
			return fmt.Errorf("installing spanner omni: %w", err)
		}
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
	log.Printf("--- Installing MySQL %s ---", mysqlVersion)

	if commandExists("mysqld") {
		out, err := runOutput("mysqld", "--version")
		if err == nil {
			version := strings.TrimSpace(out)
			log.Printf("mysql is already installed: %s", version)
			if isMySQLVersionOK(version) {
				log.Printf("mysql version is %d+, skipping installation", mysqlMajor)
				return nil
			}
			log.Printf("mysql version is too old, upgrading to MySQL %s", mysqlVersion)
			// Stop existing MySQL before upgrading
			_ = exec.Command("sudo", "service", "mysql", "stop").Run()
			_ = exec.Command("sudo", "pkill", "-f", "mysqld").Run()
			time.Sleep(2 * time.Second)
			// Remove old MySQL packages to avoid conflicts
			log.Println("removing old mysql packages")
			_ = run("sudo", "apt-get", "remove", "-y", "mysql-server", "mysql-client", "mysql-common",
				"mysql-server-core-*", "mysql-client-core-*")
			// Clear old data directory so the new release can initialize fresh
			log.Println("clearing old mysql data directory")
			_ = run("sudo", "rm", "-rf", "/var/lib/mysql")
			_ = run("sudo", "mkdir", "-p", "/var/lib/mysql")
			_ = run("sudo", "chown", "mysql:mysql", "/var/lib/mysql")
		}
	}

	major, minor, _ := strings.Cut(mysqlVersion, ".")
	minor, _, _ = strings.Cut(minor, ".")
	bundleURL := fmt.Sprintf("https://dev.mysql.com/get/Downloads/MySQL-%s.%s/mysql-server_%s-1ubuntu24.04_amd64.deb-bundle.tar", major, minor, mysqlVersion)
	bundleTar := fmt.Sprintf("/tmp/mysql-server-%s-bundle.tar", mysqlVersion)
	extractDir := "/tmp/mysql-server-" + mysqlVersion

	if _, err := os.Stat(bundleTar); err != nil {
		log.Printf("downloading MySQL %s bundle from %s", mysqlVersion, bundleURL)
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

	log.Printf("mysql %s installed successfully", mysqlVersion)
	return nil
}

// ---- start ----

func runStart(services map[string]bool) error {
	log.Println("=== Starting databases ===")

	if services["postgresql"] {
		if err := startPostgreSQL(); err != nil {
			return fmt.Errorf("starting postgresql: %w", err)
		}
	}

	if services["mysql"] {
		if err := startMySQL(); err != nil {
			return fmt.Errorf("starting mysql: %w", err)
		}
	}

	if services["spanner"] {
		if err := startSpannerOmni(); err != nil {
			return fmt.Errorf("starting spanner omni: %w", err)
		}
	}

	log.Println("=== Databases are running and configured ===")
	if services["postgresql"] {
		log.Println("PostgreSQL:   postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	}
	if services["mysql"] {
		log.Println("MySQL:        root:mysecretpassword@tcp(127.0.0.1:3306)/mysql")
	}
	if services["spanner"] {
		log.Println("Spanner Omni: localhost:" + omniPort)
	}
	return nil
}

// omniDir is where the Spanner Omni release is unpacked: the bin directory
// holds the spanner launcher and spanner_server, and data holds what a
// started server writes.
func omniDir() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, "sqlc-spanner-omni", omniVersion), nil
}

func omniLauncher(dir string) string {
	return filepath.Join(dir, "google", "spanner", "bin", "spanner")
}

// installSpannerOmni downloads the Spanner Omni server release into the
// cache and unpacks it, checking the download against the pinned SHA-256.
// It is a no-op when the release is already unpacked, and skips platforms
// the server is not published for.
func installSpannerOmni() error {
	log.Printf("--- Installing Spanner Omni %s ---", omniVersion)

	platform := runtime.GOOS + "/" + runtime.GOARCH
	bin, ok := omniBinaries[platform]
	if !ok {
		log.Printf("spanner omni is not published for %s, skipping", platform)
		return nil
	}

	dir, err := omniDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(omniLauncher(dir)); err == nil {
		log.Printf("spanner omni %s is already installed in %s", omniVersion, dir)
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	archive := filepath.Join(dir, "server.tar.gz")
	log.Printf("downloading %s", bin.URL)
	if err := downloadFile(archive, bin.URL); err != nil {
		return fmt.Errorf("downloading spanner omni: %w", err)
	}
	defer os.Remove(archive)

	sum, err := sha256File(archive)
	if err != nil {
		return err
	}
	if sum != bin.SHA256 {
		return fmt.Errorf("spanner omni download has SHA-256 %s, want %s", sum, bin.SHA256)
	}

	log.Printf("unpacking into %s", dir)
	if err := run("tar", "-xzf", archive, "-C", dir); err != nil {
		return fmt.Errorf("unpacking spanner omni: %w", err)
	}
	if _, err := os.Stat(omniLauncher(dir)); err != nil {
		return fmt.Errorf("spanner omni release does not hold %s", omniLauncher(dir))
	}
	return nil
}

// startSpannerOmni starts a Spanner Omni single server in the background,
// serving plaintext gRPC on omniPort, and waits until it accepts
// connections. It is a no-op when a server is already listening, and skips
// platforms the server is not installed on.
func startSpannerOmni() error {
	log.Println("--- Starting Spanner Omni ---")

	if omniReady() {
		log.Println("spanner omni is already running and accepting connections")
		return nil
	}

	dir, err := omniDir()
	if err != nil {
		return err
	}
	launcher := omniLauncher(dir)
	if _, err := os.Stat(launcher); err != nil {
		if _, ok := omniBinaries[runtime.GOOS+"/"+runtime.GOARCH]; !ok {
			log.Printf("spanner omni is not published for %s/%s, skipping", runtime.GOOS, runtime.GOARCH)
			return nil
		}
		return fmt.Errorf("spanner omni is not installed: run `sqlc-test-setup install` first")
	}

	data := filepath.Join(dir, "data")
	if err := os.MkdirAll(data, 0o755); err != nil {
		return err
	}
	logFile, err := os.Create(filepath.Join(dir, "server.log"))
	if err != nil {
		return err
	}
	defer logFile.Close()

	// The launcher supervises the server processes for as long as it runs,
	// so it is detached from this process and left running.
	cmd := exec.Command(launcher, "start-single-server", "--base-dir", data)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = detachedProcess()
	log.Printf("starting %s start-single-server --base-dir %s", launcher, data)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting spanner omni: %w", err)
	}
	if err := cmd.Process.Release(); err != nil {
		return err
	}

	log.Println("waiting for spanner omni to accept connections")
	if err := waitForSpannerOmni(3 * time.Minute); err != nil {
		return fmt.Errorf("spanner omni did not start in time (see %s): %w", logFile.Name(), err)
	}
	log.Println("spanner omni is accepting connections")
	return nil
}

// omniReady reports whether something accepts connections on the Spanner
// Omni gRPC port.
func omniReady() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+omniPort, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// waitForSpannerOmni polls until the server accepts connections or the
// timeout expires.
func waitForSpannerOmni(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if omniReady() {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("timed out after %s waiting for spanner omni", timeout)
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
