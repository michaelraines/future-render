package backend

import (
	"fmt"
	"sync"
)

// Factory creates a new Device instance for a named backend.
type Factory func() Device

var (
	registryMu sync.RWMutex
	registry   = make(map[string]Factory)
)

// Register registers a backend factory under the given name.
// Typically called from init() in backend implementation packages.
func Register(name string, factory Factory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, dup := registry[name]; dup {
		panic(fmt.Sprintf("backend: Register called twice for %q", name))
	}
	registry[name] = factory
}

// Create creates a new Device for the named backend. Returns an error if
// the backend is not registered.
func Create(name string) (Device, error) {
	registryMu.RLock()
	factory, ok := registry[name]
	registryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("backend: unknown backend %q", name)
	}
	return factory(), nil
}

// Available returns the names of all registered backends.
func Available() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

// IsRegistered reports whether a backend with the given name is registered.
func IsRegistered(name string) bool {
	registryMu.RLock()
	defer registryMu.RUnlock()
	_, ok := registry[name]
	return ok
}

// Resolve selects a backend by name. If name is "auto", it tries each entry
// in preferred order and returns the first registered backend. If none of the
// preferred backends are registered, it returns an error. For any other name
// it delegates to Create.
func Resolve(name string, preferred []string) (Device, string, error) {
	if name != "auto" {
		dev, err := Create(name)
		return dev, name, err
	}
	for _, p := range preferred {
		if IsRegistered(p) {
			dev, err := Create(p)
			if err != nil {
				return nil, "", err
			}
			return dev, p, nil
		}
	}
	return nil, "", fmt.Errorf("backend: no preferred backend available (tried %v, registered: %v)", preferred, Available())
}

// resetRegistry clears all registered backends. For testing only.
func resetRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = make(map[string]Factory)
}
