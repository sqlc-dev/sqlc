// Package sqlcdebug parses the SQLCDEBUG environment variable and exposes
// its key=value settings to the rest of the codebase.
//
// The SQLCDEBUG variable is a comma-separated list of name=value pairs:
//
//	SQLCDEBUG=dumpast=1,trace=trace.out
//
// Settings are looked up at the call site by declaring a package-level
// variable:
//
//	var dumpAST = sqlcdebug.New("dumpast")
//
//	func parse() {
//	    if dumpAST.Value() == "1" { ... }
//	}
//
// New panics if name is not registered in [Settings] below. Adding a new
// setting therefore requires extending the [Settings] table so that all
// known keys are documented in one place.
//
// This package is modeled after Go's internal/godebug. Unlike Go, sqlc is
// short-lived, so settings are parsed once at process startup; the
// [Update] hook exists for tests that need to reparse mid-run.
package sqlcdebug

import (
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

// Info documents a single SQLCDEBUG setting.
type Info struct {
	// Name is the SQLCDEBUG key, e.g. "dumpast".
	Name string
	// Description is a short human-readable description.
	Description string
	// Default is the value returned by [Setting.Value] when the key is
	// not present in SQLCDEBUG.
	Default string
}

// Settings is the master table of all known SQLCDEBUG keys. New keys must
// be added here before they can be looked up via [New].
var Settings = []Info{
	{Name: "dumpast", Description: "print the AST of every SQL statement"},
	{Name: "dumpcatalog", Description: "print the parsed database schema"},
	{Name: "trace", Description: "write a runtime trace to the named file (1 means trace.out)"},
	{Name: "processplugins", Description: "set to 0 to disable process-based plugins", Default: "1"},
	{Name: "databases", Description: "set to 'managed' to disable database connections via URI"},
	{Name: "dumpvetenv", Description: "print the variables available to a vet rule during evaluation"},
	{Name: "dumpexplain", Description: "print the JSON-formatted output from EXPLAIN during vet evaluation"},
}

func info(name string) (Info, bool) {
	for _, s := range Settings {
		if s.Name == name {
			return s, true
		}
	}
	return Info{}, false
}

// Setting is a single SQLCDEBUG key. Obtain one with [New]; reading it
// with [Setting.Value] returns the value parsed from the SQLCDEBUG
// environment variable, or the registered default.
type Setting struct {
	info  Info
	set   atomic.Bool
	value atomic.Pointer[string]
}

var (
	registryMu sync.Mutex
	registry   = map[string]*Setting{}
)

// New returns the Setting for name. The same pointer is returned for
// repeated calls. New panics if name is not present in [Settings].
func New(name string) *Setting {
	registryMu.Lock()
	defer registryMu.Unlock()
	if s, ok := registry[name]; ok {
		return s
	}
	i, ok := info(name)
	if !ok {
		panic("sqlcdebug: unknown setting " + name)
	}
	s := &Setting{info: i}
	registry[name] = s
	apply(s, parsedEnv())
	return s
}

// Name returns the setting's key.
func (s *Setting) Name() string { return s.info.Name }

// Default returns the value used when the key is absent from SQLCDEBUG.
func (s *Setting) Default() string { return s.info.Default }

// Value returns the parsed value of the SQLCDEBUG setting, falling back
// to the registered default.
func (s *Setting) Value() string {
	if v := s.value.Load(); v != nil {
		return *v
	}
	return s.info.Default
}

// IsSet reports whether the key was present in SQLCDEBUG.
func (s *Setting) IsSet() bool { return s.set.Load() }

// Any reports whether SQLCDEBUG contained any recognized key=value pair.
// It mirrors the legacy "is debug active?" check.
func Any() bool {
	for _, s := range registry {
		if s.IsSet() {
			return true
		}
	}
	return parsedEnvHasAny()
}

var (
	envOnce sync.Once
	envMap  map[string]string
)

func parsedEnv() map[string]string {
	envOnce.Do(func() {
		envMap = parse(os.Getenv("SQLCDEBUG"))
	})
	return envMap
}

func parsedEnvHasAny() bool {
	return len(parsedEnv()) > 0
}

func parse(raw string) map[string]string {
	out := map[string]string{}
	if raw == "" {
		return out
	}
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		k, v, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		out[k] = v
	}
	return out
}

func apply(s *Setting, env map[string]string) {
	if v, ok := env[s.info.Name]; ok {
		s.value.Store(&v)
		s.set.Store(true)
	} else {
		s.value.Store(nil)
		s.set.Store(false)
	}
}

// Update reparses the given SQLCDEBUG-formatted string and refreshes
// every registered setting. Intended for tests; production code should
// rely on the value parsed at startup.
func Update(raw string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	envMap = parse(raw)
	envOnce.Do(func() {}) // mark as parsed so future New() calls see envMap.
	for _, s := range registry {
		apply(s, envMap)
	}
}
