package plugin

import (
	"fmt"
	"sync"
)

// Registry manages the collection of available parsers
type Registry struct {
	parsers map[string]Parser
	mu      sync.RWMutex
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]Parser),
	}
}

// Register adds a new parser to the registry
func (r *Registry) Register(parser Parser) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := parser.Name()
	if _, exists := r.parsers[name]; exists {
		return fmt.Errorf("parser %s already registered", name)
	}

	r.parsers[name] = parser
	return nil
}

// GetParser retrieves a parser by name
func (r *Registry) GetParser(name string) (Parser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parser, exists := r.parsers[name]
	if !exists {
		return nil, fmt.Errorf("parser %s not found", name)
	}

	return parser, nil
}

// ListParsers returns a list of all registered parser names
func (r *Registry) ListParsers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.parsers))
	for name := range r.parsers {
		names = append(names, name)
	}
	return names
} 