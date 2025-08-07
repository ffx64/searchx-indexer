package multidb

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type multidb struct {
	dbs map[string]*sql.DB
	mu  sync.RWMutex
}

func New() *multidb {
	return &multidb{
		dbs: make(map[string]*sql.DB),
	}
}

func (m *multidb) ConnectDatabase(name, dsn string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("could not open database %s: %w", name, err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("could not ping database %s: %w", name, err)
	}

	m.dbs[name] = db
	return nil
}

func (m *multidb) GetConnection(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, exists := m.dbs[name]
	if !exists {
		return nil, fmt.Errorf("database %s not found", name)
	}
	return db, nil
}

func (m *multidb) CloseAllConnections() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for name, db := range m.dbs {
		if err := db.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close database %s: %w", name, err)
		}
	}
	return firstErr
}
