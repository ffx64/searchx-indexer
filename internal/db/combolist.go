package db

import (
	"errors"
	"fmt"
)

// ModelComboListFile represents a file combolist in the database, containing metadata about a file.
type ModelComboListFile struct {
	ID                    int64  // ID is the unique identifier for the file combolist.
	Name                  string // Name is the name of the file.
	Size                  int64  // Size is the size of the file in bytes.
	Hash                  string // Hash is the hash of the file for integrity checks.
	Status                string // Status represents the current status of the file combolist.
	Source                string // Source indicates the origin of the file.
	Type                  string // Type represents the file's type.
	Description           string // Description provides a brief description of the file.
	ProcessedEntriesCount int64  // ProcessedEntriesCount is the number of processed entries in the file.
	CreatedAt             string // CreatedAt is the timestamp when the file combolist was created.
	ProcessedAt           string // ProcessedAt is the timestamp when the file combolist was processed.
}

// ModelComboListEntrie represents an entry in the combolist file, containing sensitive data.
type ModelComboListEntrie struct {
	URL         string // URL is the URL associated with the combolist entry.
	Username    string // Username is the username for the entry.
	Password    string // Password is the password associated with the entry.
	FileLine    int64  // FileLine is the line number in the combolist file where this entry was found.
	CreatedAt   string // CreatedAt is the timestamp when the entry was created.
	ProcessedAt string // ProcessedAt is the timestamp when the entry was processed.
}

// DoesComboListFile checks if a file with the given hash already exists in the database.
//
// Parameters:
// - fileHash: The hash of the file to check for in the database.
//
// Returns:
// - bool: true if the file combolist exists, false otherwise.
// - error: an error if there was an issue with the query or database operation.
func (d *Database) DoesComboListFile(fileHash string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM file_combolist WHERE file_hash = $1)`

	err := d.cursor.QueryRow(query, fileHash).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error checking if file combolist exists: %v", err)
	}

	return exists, nil
}

// CreateComboListFile creates a new file combolist in the database if it does not already exist.
//
// Parameters:
// - file: A ModelComboListFile struct containing the file details to be inserted.
//
// Returns:
// - int64: The ID of the newly inserted file combolist.
// - error: An error if there was an issue with the query or the insertion process.
func (d *Database) CreateComboListFile(file ModelComboListFile) (int64, error) {
	exists, err := d.DoesComboListFile(file.Hash)

	if err != nil {
		return 0, err
	}

	if exists {
		return 0, errors.New("file combolist already exists, skipping create")
	}

	query := `INSERT INTO file_combolist (file_name, file_size, file_hash, status, source, file_type, file_description, processed_entries_count, created_at, processed_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	var id int64

	if err := d.cursor.QueryRow(query, file.Name, file.Size, file.Hash, file.Status, file.Source, file.Type, file.Description, file.ProcessedEntriesCount, file.CreatedAt, file.ProcessedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("error inserting file combolist: %v", err)
	}

	return id, nil
}

// GetComboListFileByID retrieves a file combolist from the database by its ID.
//
// Parameters:
// - id: The unique ID of the file combolist.
//
// Returns:
// - *ModelComboListFile: A pointer to the file combolist if found, or nil if not found.
// - error: An error if there was an issue fetching the file combolist from the database.
func (d *Database) GetComboListFileByID(id int64) (*ModelComboListFile, error) {
	file := &ModelComboListFile{}

	query := `SELECT id, file_name, file_size, file_hash, status, source, file_type, file_description, processed_entries_count, created_at, processed_at FROM file_combolist WHERE id = $1`

	err := d.cursor.QueryRow(query, id).Scan(
		&file.ID, &file.Name, &file.Size, &file.Hash, &file.Status, &file.Source, &file.Type,
		&file.Description, &file.ProcessedEntriesCount, &file.CreatedAt, &file.ProcessedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error fetching file combolist: %v", err)
	}

	return file, nil
}

// UpdateComboListFile updates an existing file combolist in the database.
//
// Parameters:
// - file: A ModelComboListFile struct containing the updated details for the file combolist.
//
// Returns:
// - error: An error if there was an issue with the update operation.
func (d *Database) UpdateComboListFile(file ModelComboListFile) error {
	query := `
		UPDATE file_combolist SET file_name = $1, file_size = $2, status = $3, source = $4, file_type = $5, 
		file_description = $6, processed_entries_count = $7, processed_at = $8 WHERE id = $9
	`
	_, err := d.cursor.Exec(query, file.Name, file.Size, file.Status, file.Source, file.Type, file.Description,
		file.ProcessedEntriesCount, file.ProcessedAt, file.ID)

	if err != nil {
		return fmt.Errorf("error updating file combolist: %v", err)
	}

	return nil
}

// DeleteComboListFile deletes a file combolist from the database by its ID.
//
// Parameters:
// - id: The unique ID of the file combolist to be deleted.
//
// Returns:
// - error: An error if there was an issue with the deletion process.
func (d *Database) DeleteComboListFile(id int64) error {
	query := `DELETE FROM file_combolist WHERE id = $1`

	_, err := d.cursor.Exec(query, id)

	if err != nil {
		return fmt.Errorf("error deleting file combolist: %v", err)
	}

	return nil
}

// DoesComboListEntrie checks if a combolist entry with the given details already exists in the database.
//
// Parameters:
// - url: The URL associated with the combolist entry.
// - username: The username associated with the combolist entry.
// - password: The password associated with the combolist entry.
//
// Returns:
// - bool: true if the combolist entry exists, false otherwise.
// - error: An error if there was an issue with the query or database operation.
func (d *Database) DoesComboListEntrie(url string, username string, password string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM combolist_entries WHERE TRIM(url) = TRIM($1) AND TRIM(username) ILIKE TRIM($2) AND TRIM(password) ILIKE TRIM($3))`

	if err := d.cursor.QueryRow(query, url, username, password).Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if combolist entry exists: %v", err)
	}

	return exists, nil
}

// CreateComboListEntrie creates a new combolist entry in the database if it does not already exist.
//
// Parameters:
// - entry: A ModelComboListEntrie struct containing the combolist entry details to be inserted.
// - fileLogID: The ID of the related file combolist.
//
// Returns:
// - int64: The ID of the newly inserted combolist entry.
// - error: An error if there was an issue with the query or the insertion process.
func (d *Database) CreateComboListEntrie(entry ModelComboListEntrie, fileLogID int64) (int64, error) {
	query := `INSERT INTO combolist_entries (file_combolist_id, url, username, password, file_line, created_at, processed_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var id int64

	if err := d.cursor.QueryRow(query, fileLogID, entry.URL, entry.Username, entry.Password, entry.FileLine, entry.CreatedAt, entry.ProcessedAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("error inserting combolist entry: %v", err)
	}

	return id, nil
}

// GetComboListEntrie retrieves a combolist entry from the database by its associated file combolist ID and entry details.
//
// Parameters:
// - fileLogID: The ID of the related file combolist.
// - url: The URL associated with the combolist entry.
// - username: The username associated with the combolist entry.
// - password: The password associated with the combolist entry.
//
// Returns:
// - *ModelComboListEntrie: A pointer to the combolist entry if found, or nil if not found.
// - error: An error if there was an issue fetching the combolist entry.
func (d *Database) GetComboListEntrie(fileLogID int64, url string, username string, password string) (*ModelComboListEntrie, error) {
	entry := &ModelComboListEntrie{}

	query := `SELECT url, username, password, created_at, processed_at FROM combolist_entries WHERE file_combolist_id = $1 AND url = $2 AND username = $3 AND password = $4`

	if err := d.cursor.QueryRow(query, fileLogID, url, username, password).Scan(&entry.URL, &entry.Username, &entry.Password, &entry.CreatedAt, &entry.ProcessedAt); err != nil {
		return nil, fmt.Errorf("error fetching combolist entry: %v", err)
	}

	return entry, nil
}

// UpdateComboListEntrie updates an existing combolist entry in the database.
//
// Parameters:
// - entry: A ModelComboListEntrie struct containing the updated details for the combolist entry.
// - fileLogID: The ID of the related file combolist.
//
// Returns:
// - int64: The ID of the updated combolist entry.
// - error: An error if there was an issue with the update operation.
func (d *Database) UpdateComboListEntrie(entry ModelComboListEntrie, fileLogID int64) (int64, error) {
	query := `UPDATE combolist_entries SET url = $1, username = $2, password = $3, processed_at = $4 WHERE file_combolist_id = $5 AND url = $6 AND username = $7 RETURNING id`

	var id int64

	if err := d.cursor.QueryRow(query, entry.URL, entry.Username, entry.Password, entry.ProcessedAt, fileLogID, entry.URL, entry.Username).Scan(&id); err != nil {
		return 0, fmt.Errorf("error updating combolist entry: %v", err)
	}

	return id, nil
}

// DeleteComboListEntrie deletes a combolist entry from the database based on the file combolist ID and entry details.
//
// Parameters:
// - fileLogID: The ID of the related file combolist.
// - url: The URL associated with the combolist entry.
// - username: The username associated with the combolist entry.
//
// Returns:
// - error: An error if there was an issue with the deletion process.
func (d *Database) DeleteComboListEntrie(fileLogID int64, url string, username string) error {
	query := `DELETE FROM combolist_entries WHERE file_combolist_id = $1 AND url = $2 AND username = $3 RETURNING id`

	var id int64

	if err := d.cursor.QueryRow(query, fileLogID, url, username).Scan(&id); err != nil {
		return fmt.Errorf("error deleting combolist entry: %v", err)
	}

	return nil
}
