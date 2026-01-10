package engine

import (
	"sync"
)

var registerOnce sync.Once

// RegisterBuiltinEngines registers all built-in database engines.
// This function should be called once during application initialization.
// It is safe to call multiple times - subsequent calls are no-ops.
func RegisterBuiltinEngines(factories map[string]EngineFactory) {
	registerOnce.Do(func() {
		for name, factory := range factories {
			Register(name, factory)
		}
	})
}
