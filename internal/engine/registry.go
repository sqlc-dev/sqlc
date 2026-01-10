package engine

import (
	"fmt"
	"sync"
)

// Registry is a global registry of database engines.
// It allows both built-in and plugin engines to be registered and retrieved.
type Registry struct {
	mu      sync.RWMutex
	engines map[string]EngineFactory
}

// globalRegistry is the default engine registry used by the application.
var globalRegistry = &Registry{
	engines: make(map[string]EngineFactory),
}

// Register adds a new engine factory to the global registry.
// It panics if an engine with the same name is already registered.
func Register(name string, factory EngineFactory) {
	globalRegistry.Register(name, factory)
}

// Get retrieves an engine by name from the global registry.
// It returns an error if the engine is not found.
func Get(name string) (Engine, error) {
	return globalRegistry.Get(name)
}

// List returns a list of all registered engine names.
func List() []string {
	return globalRegistry.List()
}

// IsRegistered returns true if an engine with the given name is registered.
func IsRegistered(name string) bool {
	return globalRegistry.IsRegistered(name)
}

// Register adds a new engine factory to this registry.
// It panics if an engine with the same name is already registered.
func (r *Registry) Register(name string, factory EngineFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.engines[name]; exists {
		panic(fmt.Sprintf("engine %q is already registered", name))
	}
	r.engines[name] = factory
}

// RegisterOrReplace adds or replaces an engine factory in this registry.
// This is useful for testing or for replacing built-in engines with plugins.
func (r *Registry) RegisterOrReplace(name string, factory EngineFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.engines[name] = factory
}

// Get retrieves an engine by name from this registry.
// It returns an error if the engine is not found.
func (r *Registry) Get(name string) (Engine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.engines[name]
	if !ok {
		return nil, fmt.Errorf("unknown engine: %s", name)
	}
	return factory(), nil
}

// List returns a list of all registered engine names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.engines))
	for name := range r.engines {
		names = append(names, name)
	}
	return names
}

// IsRegistered returns true if an engine with the given name is registered.
func (r *Registry) IsRegistered(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.engines[name]
	return ok
}

// Unregister removes an engine from this registry.
// This is primarily useful for testing.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.engines, name)
}
