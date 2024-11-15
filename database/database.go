package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(sqliteFile string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", sqliteFile)

	if err != nil {
		return nil, fmt.Errorf("[!] error connecting to the database: %v", err)
	}
	print("[+] database connected successfully!\n")
	return conn, nil
}

func CredentialsCreateTable(conn *sql.DB) error {
	query := `
	create table if not exists credentials (
		id integer primary key autoincrement,
		url text not null,
		user text not null,
		pass text,
		file_name text not null,
		file_size integer not null,
		description text,
		created_at datetime default current_timestamp,
		unique(url, user, pass)
	);`

	_, err := conn.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func CredentialsInsertTable(conn *sql.DB, url string, username string, password string, description string, filename string, filesize int64) error {
	query := `
	insert or ignore into credentials (
		url,
		user,
		pass,
		description,
		file_name,
		file_size
	) values (?, ?, ?, ?, ?, ?);`

	tx, _ := conn.Begin()

	_, err := conn.Exec(query, url, username, password, description, filename, filesize)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func GoroutineCredentialsInsertTable(conn *sql.DB, url string, username string, password string, description string, filename string, filesize int64) {
	err := CredentialsInsertTable(conn, url, username, password, description, filename, filesize)

	if err != nil {
		print("[!] error credentials insert: ", err)
	}
}

func DisableSecuritySettings(conn *sql.DB) error {
	_, err := conn.Exec("PRAGMA synchronous = OFF;")
	if err != nil {
		return err
	}
	print("[+] synchronous mode disable\n")

	_, err = conn.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return err
	}
	print("[+] journal mode set to WAL\n")

	return nil
}

func SetTempStoreMemorySpeed(conn *sql.DB) error {
	_, err := conn.Exec("PRAGMA temp_store = MEMORY;")

	if err != nil {
		return err
	}
	print("[+] temp store = memory\n")

	return nil
}
