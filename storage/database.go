package storage

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq" // Importa driver PostgreSQL
)

type DBManager struct {
	dbs map[string]*sql.DB
	mu  sync.RWMutex
}

func NewDBManager() *DBManager {
	return &DBManager{
		dbs: make(map[string]*sql.DB),
	}
}

// AddDB adiciona uma nova conexão de banco ao gerenciador
func (m *DBManager) AddDB(name string, dsn string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com banco %s: %w", name, err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar com o banco %s: %w", name, err)
	}

	m.dbs[name] = db
	return nil
}

// GetDB retorna a instância de conexão com base no nome
func (m *DBManager) GetDB(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, exists := m.dbs[name]
	if !exists {
		return nil, fmt.Errorf("banco %s não encontrado", name)
	}
	return db, nil
}

// CloseAll fecha todas as conexões abertas
func (m *DBManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error
	for name, db := range m.dbs {
		if err := db.Close(); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("erro ao fechar banco %s: %w", name, err)
			}
		}
	}
	return firstErr
}
