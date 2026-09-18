package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// mssqlRelease is the SQL Server release installed from Microsoft's
	// apt repository, the one internal/goldeneye generates the dialect
	// from.
	mssqlRelease = "2025"

	// mssqlSAPassword is the sa password the server is set up with: the
	// one docker-compose.yml and CI use.
	mssqlSAPassword = "Mysecretpassword1!"

	mssqlPort   = "1433"
	mssqlServer = "/opt/mssql/bin/sqlservr"
	mssqlConf   = "/opt/mssql/bin/mssql-conf"
	mssqlLog    = "/var/opt/mssql/log/errorlog"
)

// mssqlUbuntuReleases lists the Ubuntu releases Microsoft publishes SQL
// Server packages for.
var mssqlUbuntuReleases = map[string]bool{"22.04": true, "24.04": true}

// installMSSQL installs SQL Server from Microsoft's apt repository and
// runs its setup non-interactively, accepting the EULA on the user's
// behalf. It is a no-op when the server is already installed, and skips
// platforms the packages are not published for: Ubuntu 22.04 and 24.04 on
// amd64 and arm64.
func installMSSQL() error {
	log.Printf("--- Installing SQL Server %s ---", mssqlRelease)

	if _, err := os.Stat(mssqlServer); err == nil {
		log.Printf("sql server is already installed at %s", mssqlServer)
		return nil
	}

	ubuntu, err := ubuntuRelease()
	if err != nil {
		log.Printf("sql server packages are published for Ubuntu only, skipping: %s", err)
		return nil
	}
	if runtime.GOOS != "linux" || (runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64") || !mssqlUbuntuReleases[ubuntu] {
		log.Printf("sql server packages are not published for Ubuntu %s on %s/%s, skipping", ubuntu, runtime.GOOS, runtime.GOARCH)
		return nil
	}

	tmp, err := os.MkdirTemp("", "sqlc-mssql-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	key := filepath.Join(tmp, "microsoft.asc")
	log.Println("adding microsoft's package signing key")
	if err := downloadFile(key, "https://packages.microsoft.com/keys/microsoft.asc"); err != nil {
		return fmt.Errorf("downloading the signing key: %w", err)
	}
	if err := run("sudo", "gpg", "--dearmor", "--yes", "-o", "/usr/share/keyrings/microsoft-prod.gpg", key); err != nil {
		return fmt.Errorf("installing the signing key: %w", err)
	}

	list := filepath.Join(tmp, "mssql-server.list")
	listURL := "https://packages.microsoft.com/config/ubuntu/" + ubuntu + "/mssql-server-" + mssqlRelease + ".list"
	log.Printf("adding the apt repository from %s", listURL)
	if err := downloadFile(list, listURL); err != nil {
		return fmt.Errorf("downloading the repository list: %w", err)
	}
	if err := run("sudo", "cp", list, "/etc/apt/sources.list.d/mssql-server-"+mssqlRelease+".list"); err != nil {
		return err
	}

	if err := run("sudo", "apt-get", "update"); err != nil {
		return fmt.Errorf("apt-get update: %w", err)
	}
	log.Println("installing mssql-server")
	if err := run("sudo", "env", "DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y", "mssql-server"); err != nil {
		return fmt.Errorf("apt-get install mssql-server: %w", err)
	}

	// Setup writes the server's configuration and accepts the EULA, then
	// starts the service through systemd, which fails where there is none.
	// The configuration is what install needs; start brings the server up.
	log.Println("running mssql-conf setup")
	err = run("sudo", "env", "ACCEPT_EULA=Y", "MSSQL_PID=Developer", "MSSQL_SA_PASSWORD="+mssqlSAPassword,
		mssqlConf, "-n", "setup", "accept-eula")
	if _, statErr := os.Stat("/var/opt/mssql/mssql.conf"); statErr != nil {
		if err != nil {
			return fmt.Errorf("mssql-conf setup: %w", err)
		}
		return fmt.Errorf("mssql-conf setup wrote no /var/opt/mssql/mssql.conf")
	}
	if err != nil {
		log.Printf("mssql-conf setup wrote the configuration but could not start the service (%s); start will", err)
	}
	return nil
}

// ubuntuRelease reads the Ubuntu release from /etc/os-release.
func ubuntuRelease() (string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}
	var id, version string
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		value = strings.Trim(value, `"`)
		switch key {
		case "ID":
			id = value
		case "VERSION_ID":
			version = value
		}
	}
	if id != "ubuntu" {
		return "", fmt.Errorf("this is %s, not ubuntu", id)
	}
	return version, nil
}

// startMSSQL starts SQL Server and waits until it is ready for client
// connections: through systemd where it runs, otherwise by running the
// server in the background as its own user. It is a no-op when the server
// is already listening, and skips machines it is not installed on.
func startMSSQL() error {
	log.Println("--- Starting SQL Server ---")

	if mssqlListening() {
		log.Println("sql server is already running and accepting connections")
		return nil
	}
	if _, err := os.Stat(mssqlServer); err != nil {
		if _, err := ubuntuRelease(); err != nil || runtime.GOOS != "linux" {
			log.Println("sql server is not installed on this platform, skipping")
			return nil
		}
		return fmt.Errorf("sql server is not installed: run `sqlc-test-setup install mssql` first")
	}

	if systemdRunning() {
		log.Println("starting mssql-server through systemd")
		if err := run("sudo", "systemctl", "start", "mssql-server"); err != nil {
			return fmt.Errorf("systemctl start mssql-server: %w", err)
		}
	} else {
		// The server's own log is mssqlLog; what it prints goes to a file
		// of its own rather than this process's output, which a detached
		// process would otherwise hold open.
		cache, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		dir := filepath.Join(cache, "sqlc-mssql")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		logFile, err := os.Create(filepath.Join(dir, "sqlservr.log"))
		if err != nil {
			return err
		}
		defer logFile.Close()

		log.Printf("starting %s in the background", mssqlServer)
		cmd := exec.Command("sudo", "-u", "mssql", "env", "ACCEPT_EULA=Y", "MSSQL_SA_PASSWORD="+mssqlSAPassword, mssqlServer)
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		cmd.SysProcAttr = detachedProcess()
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("starting sqlservr: %w", err)
		}
		if err := cmd.Process.Release(); err != nil {
			return err
		}
	}

	log.Println("waiting for sql server to accept connections")
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		if mssqlListening() && mssqlReadyForClients() {
			log.Println("sql server is accepting connections")
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("sql server did not start in time (see %s)", mssqlLog)
}

// systemdRunning reports whether systemd is the init system here, so a
// service can be started through it.
func systemdRunning() bool {
	out, err := exec.Command("systemctl", "is-system-running").Output()
	state := strings.TrimSpace(string(out))
	if err != nil && state == "" {
		return false
	}
	return state != "offline" && state != "unknown"
}

// mssqlListening reports whether something accepts connections on the SQL
// Server port.
func mssqlListening() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+mssqlPort, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// mssqlReadyForClients reports whether the server's log says it is ready
// for client connections, which comes a moment after it starts listening.
// The log is readable by the mssql user only, so it is read through sudo.
func mssqlReadyForClients() bool {
	return exec.Command("sudo", "grep", "-q", "ready for client connections", mssqlLog).Run() == nil
}
