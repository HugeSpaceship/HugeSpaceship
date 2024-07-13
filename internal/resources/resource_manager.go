package resources

import (
	"cmp"
	"errors"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"io"
	"log/slog"
	"slices"
)

type ResourceManager struct {
	priorities  []connectionPriority
	backends    map[string]backends.ResourceBackend
	connections map[string]backends.BackendConnection
	config      *config.Config
}

type connectionPriority struct {
	name     string
	priority uint
}

func NewResourceManager(cfg *config.Config) *ResourceManager {
	return &ResourceManager{
		backends:    map[string]backends.ResourceBackend{},
		connections: map[string]backends.BackendConnection{},
		priorities:  []connectionPriority{},
		config:      cfg,
	}
}

func (r *ResourceManager) RegisterBackend(name string, backend backends.ResourceBackend) {
	r.backends[name] = backend
}

func (r *ResourceManager) RegisterBackendConfig(cfg *config.ResourceBackendConfig) error {
	backend, ok := r.backends[cfg.Type]
	if !ok {
		return fmt.Errorf("unknown backend: %s", cfg.Type)
	}

	conn, err := backend.InitConnection(cfg.Config, r.config)
	if err != nil {
		return err
	}
	r.connections[cfg.Name] = conn
	r.priorities = append(r.priorities, connectionPriority{priority: cfg.Priority, name: cfg.Name})

	slices.SortStableFunc(r.priorities, func(a, b connectionPriority) int {
		return cmp.Compare(a.priority, b.priority)
	})
	return nil
}

func (r *ResourceManager) HasResource(hash string) (bool, string, error) {
	var failError error
	for _, priority := range r.priorities {
		exists, err := r.connections[priority.name].HasResource(hash)
		if err != nil {
			slog.Error("failed to check resource", slog.String("backend", priority.name), slog.String("hash", hash))
			failError = errors.Join(failError, err)
			continue
		}
		if exists {
			return true, priority.name, failError
		}
	}
	return false, "", failError
}

func (r *ResourceManager) GetResource(hash string) (io.ReadCloser, int64, bool, error) {
	exists, backend, err := r.HasResource(hash)
	if err != nil && !exists {
		return nil, 0, false, err
	}
	if !exists {
		return nil, 0, false, nil
	}

	reader, length, err := r.connections[backend].GetResource(hash)
	return reader, length, true, err
}
