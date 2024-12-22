package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// database represents a connection to a PostgreSQL database.
type Database struct {
	host     string
	port     int32
	username string
	password string
	database string
	sslmode  string
	cursor   *sql.DB
}

// NewDatabase creates and returns a new Database instance with the provided connection parameters.
// The default SSL mode is set to "disable".
//
// Parameters:
// - host: The host where the database is located (e.g., "localhost").
// - port: The port on which the database is listening (e.g., "5432").
// - username: The username used to authenticate with the database.
// - password: The password associated with the username.
// - db: The name of the database to connect to.
//
// Returns:
// - *Database: A pointer to the newly created Database instance.
// - The function sets the default SSL mode to "disable" and configures the connection.
func NewDatabase(port int32, host, username, password, db, sslmode string) *Database {
	model := new(Database)

	model.host = host
	model.port = port
	model.username = username
	model.password = password
	model.database = db

	return model
}

// Connect establishes a connection to the PostgreSQL database and stores the connection
// in the cursor field. Returns the instance of the Database struct and any error encountered.
//
// Returns:
// - *Database: The Database instance with an open database connection.
// - error: An error object if the connection fails, otherwise nil.
func (d *Database) Connect() (*Database, error) {
	conn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.host, d.port, d.username, d.password, d.database, d.sslmode))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	d.cursor = conn

	log.Println("database connected")

	return d, nil
}

// Close closes the connection to the database. It returns an error if the closing fails.
//
// Returns:
// - error: Returns nil if the connection is closed successfully, otherwise an error.
func (d *Database) Close() error {
	if d.cursor != nil {
		return d.cursor.Close()
	}
	return fmt.Errorf("no active database connection to close")
}
