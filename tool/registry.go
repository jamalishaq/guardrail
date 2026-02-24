// tool/registry.go
package tool

import (
	"errors"
	"sync"
)

var (
	// ErrToolExists is returned when adding a tool with a duplicate name.
	ErrToolExists = errors.New("tool already exists")
	// ErrToolNotFound is returned when a tool name can't be found in the registry.
	ErrToolNotFound = errors.New("tool not found")
)

// Registry stores tools by name. Safe for concurrent access.
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Add registers a tool. Returns ErrToolExists if the name already exists.
func (r *Registry) Add(t Tool) error {
	if t == nil {
		return errors.New("tool is nil")
	}
	name := t.Name()
	if name == "" {
		return errors.New("tool name is empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tools[name]; ok {
		return ErrToolExists
	}
	r.tools[name] = t
	return nil
}

// Get returns a tool by name.
func (r *Registry) Get(name string) (Tool, error) {
	if name == "" {
		return nil, errors.New("tool name is empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tools[name]
	if !ok {
		return nil, ErrToolNotFound
	}
	return t, nil
}

// ListNames returns all tool names (unordered).
func (r *Registry) ListNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.tools))
	for name := range r.tools {
		out = append(out, name)
	}
	return out
}

// ListTools returns all registered tools (unordered).
// Note: tools are returned as interfaces; callers should treat them as read-only.
func (r *Registry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t)
	}
	return out
}