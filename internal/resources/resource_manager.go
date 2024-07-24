// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/config"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/sqlc"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends/pg_lob"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"slices"
)

type ResourceManager struct {
	priorities  []connectionPriority
	backends    map[string]backends.ResourceBackend
	connections map[string]backends.BackendConnection
	v           *viper.Viper
	sql         *sqlc.Queries
}

type connectionPriority struct {
	name     string
	priority uint
}

func NewResourceManager(v *viper.Viper, pool *pgxpool.Pool) *ResourceManager {
	rm := ResourceManager{
		backends:    map[string]backends.ResourceBackend{},
		connections: map[string]backends.BackendConnection{},
		priorities:  []connectionPriority{},
		v:           v,
		sql:         sqlc.New(pool),
	}

	rm.RegisterBackend("pg_lob", &pg_lob.Backend{})

	return &rm
}

func (r *ResourceManager) RegisterBackend(name string, backend backends.ResourceBackend) {
	r.backends[name] = backend
}

func (r *ResourceManager) RegisterBackendConfig(cfg *config.ResourceBackendConfig) error {
	backend, ok := r.backends[cfg.Type]
	if !ok {
		return fmt.Errorf("unknown backend: %s", cfg.Type)
	}

	conn, err := backend.InitConnection(cfg.Config, r.v)
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

func (r *ResourceManager) Start() error {
	b := r.v.GetStringMap("resource-server.backends")
	for name, backend := range b {
		cfgMap := backend.(map[string]interface{})
		priority := cfgMap["priority"].(int)
		backendType := cfgMap["type"].(string)
		delete(cfgMap, "priority")
		delete(cfgMap, "type")
		backendCfg := config.ResourceBackendConfig{
			Name:     name,
			Type:     backendType,
			Priority: uint(priority),
			Config:   cfgMap,
		}
		err := r.RegisterBackendConfig(&backendCfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ResourceManager) HasResource(hash string) (bool, string, error) {
	if len(r.connections) == 0 {
		return false, "", errors.New("no resource backends have been configured")
	}

	var failError error
	for _, priority := range r.priorities {
		exists, err := r.connections[priority.name].HasResource(hash)
		if err != nil {
			slog.Error("failed to check resource", slog.String("backend", priority.name), slog.String("hash", hash), slog.Any("error", err))
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
	if len(r.connections) == 0 {
		return nil, 0, false, errors.New("no resource backends have been configured")
	}

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
	if len(r.connections) == 0 {
		return errors.New("no resource backends have been configured")
	}

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

		reader.Seek(0, io.SeekStart)
		magic := make([]byte, 4)
		reader.Read(magic)

		slog.Debug("Resource Magic number", "magic", string(magic))

		err = r.sql.InsertResource(context.Background(), sqlc.InsertResourceParams{
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
