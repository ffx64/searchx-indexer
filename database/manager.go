package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DBManager manages multiple database connections.
type DBManager struct {
	connections map[string]*sql.DB
	mu          sync.Mutex
}

var manager *DBManager
var once sync.Once

// GetManager returns a singleton instance of DBManager.
// It ensures that only one instance of the manager is created.
func GetManager() *DBManager {
	once.Do(func() {
		manager = &DBManager{
			connections: make(map[string]*sql.DB),
		}
	})
	return manager
}

// Connect establishes a new database connection and stores it in the manager.
// If a connection with the same name already exists, it returns an error.
//
// driver available : postgres and sqlite
//
// Usage:
// err := dbManager.Connect("combolist", "postgres", "localhost", 5432, "user", "password", "dbname")
//
//	if err != nil {
//	    log.Println("Failed to connect:", err)
//	}
func (m *DBManager) Connect(name, driver string, host string, port int, user, password, dbname string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.connections[name]; exists {
		return fmt.Errorf("connection '%s' already exists", name)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open(driver, dsn)

	if err != nil {
		return fmt.Errorf("error connecting to database '%s': %v", name, err)
	}

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("error validating connection to '%s': %v", name, err)
	}

	m.connections[name] = conn
	log.Printf("[+] Connected to database '%s'\n", name)
	return nil
}

// GetDB retrieves a database connection by its name.
// If the connection does not exist, it returns an error.
//
// Usage:
// db, err := dbManager.GetDB("combolist")
//
//	if err != nil {
//	    log.Println("Database not found:", err)
//	}
func (m *DBManager) GetDB(name string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("connection '%s' not found", name)
	}

	return conn, nil
}

// Close closes all database connections managed by DBManager.
//
// Usage:
// dbManager.Close()
func (m *DBManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, conn := range m.connections {
		_ = conn.Close()
		log.Printf("[-] Connection '%s' closed\n", name)
	}

	m.connections = make(map[string]*sql.DB)
}
