package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func (d *Database) NewConnection() (*Database, error) {
	defaultConfig(d)

	conn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.Username, d.Password, d.DBName, d.SSL))
	if err != nil {
		return nil, fmt.Errorf("[!] error connecting to the database: %v", err)
	}

	d.cursor = conn
	log.Println("[+] database connected successfully!")
	return d, nil
}

func (d *Database) CloseConnection() error {
	return d.cursor.Close()
}

func (d *Database) FileLogExists(fileHash string, fileName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM file_logs
			WHERE file_hash = $1 OR file_name = $2
		)
	`
	err := d.cursor.QueryRow(query, fileHash, fileName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if file log exists: %v", err)
	}
	return exists, nil
}

func (d *Database) InsertStealerFile(file StealerFileLogModel) (int64, error) {
	exists, err := d.FileLogExists(file.FileHash, file.FileName)
	if err != nil {
		return 0, err
	}
	if exists {
		log.Println("[+] file log already exists, skipping insert.")
		return 0, nil
	}

	query := `
		INSERT INTO file_logs (file_name, file_size, file_hash, status, source, file_type, file_description, processed_entries_count, created_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var fileLogID int64
	err = d.cursor.QueryRow(query, file.FileName, file.FileSize, file.FileHash, file.Status, file.Source, file.FileType, file.FileDescription, file.ProcessedEntriesCount, file.CreatedAt, file.ProcessedAt).Scan(&fileLogID)
	if err != nil {
		return 0, fmt.Errorf("error inserting file log: %v", err)
	}

	return fileLogID, nil
}

func (d *Database) LogEntryExists(fileLogID int64, url string, username string, password string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM log_entries
			WHERE file_log_id = $1 AND url = $2 AND username = $3 AND password = $4
		)
	`
	err := d.cursor.QueryRow(query, fileLogID, url, username, password).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if log entry exists: %v", err)
	}
	return exists, nil
}

func (d *Database) InsertStealerEntrie(entry StealerEntrieLogModel, fileID int64) error {
	exists, err := d.LogEntryExists(fileID, entry.URL, entry.Username, entry.Password)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("log entry already exists, skipping insert.")
	}

	query := `
		INSERT INTO log_entries (file_log_id, url, username, password, created_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = d.cursor.Exec(query, fileID, entry.URL, entry.Username, entry.Password, entry.CreatedAt, entry.ProcessedAt)
	if err != nil {
		return fmt.Errorf("error inserting log entry: %v", err)
	}

	return nil
}
