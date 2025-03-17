package combolist

import (
	"database/sql"
	"fmt"
)

// CombolistFileExists checks if a file with the given hash already exists in the database.
func CombolistFileExists(cursor *sql.DB, fileHash string) (bool, error) {
	var exists bool

	err := cursor.QueryRow(`SELECT EXISTS (SELECT 1 FROM file_combolist WHERE file_hash = $1);`, fileHash).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error checking if file combolist exists: %v", err)
	}

	return exists, nil
}

// CombolistEntryWithUrlExists checks if a combolist entry with the given details already exists in the database.
func CombolistEntryWithUrlExists(cursor *sql.DB, url string, username string, password string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM combolist_entries_user u JOIN combolist_entries_urls urls ON u.id = urls.file_combolist_id WHERE u.username = $1 AND u.password = $2 AND urls.url = $3);`

	if err := cursor.QueryRow(query, username, password, url).Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if combolist entry exists: %v", err)
	}

	return exists, nil
}

// CreateCombolistFile creates a new combolist entry in the database if it does not already exist.
func CreateCombolistFile(cursor *sql.DB, entry EntityCombolistFile) (int64, error) {
	query := `INSERT INTO file_combolist (file_name, file_size, file_hash, file_type, file_description, source, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	var id int64

	if err := cursor.QueryRow(query, entry.Name, entry.Size, entry.Hash, entry.Type, entry.Description, entry.Source, entry.CreatedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("error inserting combolist file entry: %v", err)
	}

	return id, nil
}

// CreateCombolistUserEntry creates a new combolist user entry in the database if it does not already exist.
func CreateCombolistUserEntry(cursor *sql.DB, entry EntityCombolistUserEntrie) (int64, error) {
	query := `INSERT INTO combolist_entries_user (file_combolist_id, username, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id;`

	var id int64

	if err := cursor.QueryRow(query, entry.FileID, entry.Username, entry.Password, entry.CreatedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("error inserting combolist user entry: %v", err)
	}

	return id, nil
}

// CreateCombolistUrlEntry creates a new combolist URL entry in the database if it does not already exist.
func CreateCombolistUrlEntry(cursor *sql.DB, entry EntityCombolistUrlEntrie) (int64, error) {
	query := `INSERT INTO combolist_entries_urls (combolist_entries_user_id, url, created_at) VALUES ($1, $2, $3) RETURNING id;`

	var id int64

	if err := cursor.QueryRow(query, entry.UserID, entry.URL, entry.CreatedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("error inserting combolist url entry: %v", err)
	}

	return id, nil
}

// UpdateCombolistFileProcessedEntriesCount updates the processed entries count in the database.
func UpdateCombolistFileProcessedEntriesCount(cursor *sql.DB, fileID int64, count int64) error {
	query := `UPDATE combolist_file SET processed_entries_count = $1 WHERE id = $2;`

	_, err := cursor.Exec(query, count, fileID)

	if err != nil {
		return fmt.Errorf("error updating combolist processed entries count: %v", err)
	}

	return nil
}

// MarkCombolistFileAsProcessed marks a file as processed in the database.
func MarkCombolistFileAsProcessed(cursor *sql.DB, fileID int64) error {
	query := `UPDATE combolist_file SET status = 1 WHERE id = $1;`

	_, err := cursor.Exec(query, fileID)

	if err != nil {
		return fmt.Errorf("error marking combolist file as processed: %v", err)
	}

	return nil
}

// DeleteCombolistFile deletes a combolist file entry from the database.
func DeleteCombolistFile(cursor *sql.DB, fileID int64) error {
	query := `DELETE FROM file_combolist WHERE id = $1;`

	_, err := cursor.Exec(query, fileID)
	if err != nil {
		return fmt.Errorf("error deleting combolist file: %v", err)
	}

	return nil
}

// DeleteCombolistUserEntry deletes a combolist user entry from the database.
func DeleteCombolistUserEntry(cursor *sql.DB, userID int64) error {
	query := `DELETE FROM combolist_entries_user WHERE id = $1;`

	_, err := cursor.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error deleting combolist user entry: %v", err)
	}

	return nil
}

// DeleteCombolistUrlEntry deletes a combolist URL entry from the database.
func DeleteCombolistUrlEntry(cursor *sql.DB, urlID int64) error {
	query := `DELETE FROM combolist_entries_urls WHERE id = $1;`

	_, err := cursor.Exec(query, urlID)
	if err != nil {
		return fmt.Errorf("error deleting combolist URL entry: %v", err)
	}

	return nil
}

// DeleteCombolistFileAndEntries deletes a combolist file and its associated user and URL entries.
func DeleteCombolistFileAndEntries(cursor *sql.DB, fileID int64) error {
	_, err := cursor.Exec(`DELETE FROM combolist_entries_urls WHERE combolist_entries_user_id IN (SELECT id FROM combolist_entries_user WHERE file_combolist_id = $1);`, fileID)
	if err != nil {
		return fmt.Errorf("error deleting combolist URL entries: %v", err)
	}

	_, err = cursor.Exec(`DELETE FROM combolist_entries_user WHERE file_combolist_id = $1;`, fileID)
	if err != nil {
		return fmt.Errorf("error deleting combolist user entries: %v", err)
	}

	_, err = cursor.Exec(`DELETE FROM file_combolist WHERE id = $1;`, fileID)
	if err != nil {
		return fmt.Errorf("error deleting combolist file: %v", err)
	}

	return nil
}
