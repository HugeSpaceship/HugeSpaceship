// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/sqlc"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/google/uuid"
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

func (r *ResourceManager) UploadResource(hash string, res io.Reader, length int64, user uuid.UUID) error {
	for _, priority := range r.priorities {
		if !r.connections[priority.name].CanUpload() {
			continue
		}

		buf := new(bytes.Buffer)
		io.Copy(buf, res)
		reader := bytes.NewReader(buf.Bytes())

		err := r.connections[priority.name].UploadResource(hash, reader, length)

		if err != nil {
			slog.Error("failed to upload resource", slog.String("backend", priority.name), slog.String("hash", hash))
			continue
		}

		conn, err := db.GetConnection(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get connection: %w", err)
		}

		c := sqlc.New(conn)
		reader.Seek(0, io.SeekStart)
		magic := make([]byte, 4)
		reader.Read(magic)

		fmt.Println(string(magic))

		err = c.InsertResource(context.Background(), sqlc.InsertResourceParams{
			Uploader:    user,
			Size:        length,
			Type:        sqlc.ResourceTypeLVL,
			Hash:        hash,
			Backend:     sqlc.ResourceBackendsPgLob,
			Backendname: "default",
		})

		return err
	}
	return errors.New("failed to upload resource")
}
