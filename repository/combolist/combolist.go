package combolist

import (
	"database/sql"
	"fmt"
	models "searchx-indexer/models/combolist"
)

// CombolistRepository handles all database operations related to combolist data,
// including files, user credentials, and URLs.
type CombolistRepository struct {
	db *sql.DB
}

// NewCombolistRepository initializes a new CombolistRepository with the given SQL database connection.
func NewCombolistRepository(db *sql.DB) *CombolistRepository {
	return &CombolistRepository{db: db}
}

////////////////////////////////
//           FILE             //
////////////////////////////////

// FileExists checks if a file with the given hash exists in the `file` table.
//
// Parameters:
//   - hash: SHA-256 (or similar) hash used to uniquely identify the file.
//
// Returns:
//   - true if the file exists, false otherwise.
//   - error if the query fails.
func (r *CombolistRepository) FileExists(hash string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM files WHERE hash = $1)`, hash).Scan(&exists)
	return exists, err
}

// GetFileByHash returns the ID of the file with the specified hash.
//
// Parameters:
//   - hash: The unique hash used to identify the file.
//
// Returns:
//   - ID of the file.
//   - error if the file is not found or if the query fails.
func (r *CombolistRepository) GetFileByHash(hash string) (int64, error) {
	var id int64
	err := r.db.QueryRow(`SELECT id FROM files WHERE hash = $1`, hash).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("file not found for hash: %s", hash)
	}
	return id, err
}

// CreateFile inserts a new file entry into the database, avoiding duplicates by hash.
//
// Parameters:
//   - file: File model containing metadata such as name, size, hash, etc.
//
// Returns:
//   - ID of the newly inserted file.
//   - error if the insert fails or the file already exists.
func (r *CombolistRepository) CreateFile(file models.Files) (int64, error) {
	query := `
		INSERT INTO files (
			agent_key, name, size, hash, type, description,
			source, processed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (hash) DO NOTHING
		RETURNING id
	`
	var id int64
	err := r.db.QueryRow(query,
		file.AgentKey, file.Name, file.Size, file.Hash, file.Type,
		file.Description, file.Source, file.ProcessedAt,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("file already exists")
	}

	return id, err
}

////////////////////////////////
//           USER             //
////////////////////////////////

// GetUserIDByCredentials fetches the user ID for a given username and password combination.
//
// Parameters:
//   - username: The user's login name.
//   - password: The user's password.
//
// Returns:
//   - user ID if found.
//   - 0 and no error if not found.
//   - error if query fails.
func (r *CombolistRepository) GetUserIDByCredentialsInFile(fileID int64, username, password string) (int64, error) {
	var id int64
	query := `SELECT id FROM users WHERE file_id = $1 AND username = $2 AND password = $3`
	err := r.db.QueryRow(query, fileID, username, password).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil // Not found, but not an error
	}
	return id, err
}

// CreateUserEntry inserts a new user entry linked to a specific file.
//
// Parameters:
//   - user: User model containing username, password, and file reference.
//
// Returns:
//   - ID of the new user.
//   - error if the insert fails.
func (r *CombolistRepository) CreateUserEntry(user models.Users) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (file_id, username, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRow(query, user.FileID, user.Username, user.Password).Scan(&id)
	return id, err
}

////////////////////////////////
//           URL              //
////////////////////////////////

// UrlEntryExists checks if a given URL is already associated with the specified user ID.
//
// Parameters:
//   - userID: ID of the user.
//   - url: The URL to check for existence.
//
// Returns:
//   - true if the URL exists for the user.
//   - false if not.
//   - error if the query fails.
func (r *CombolistRepository) UrlEntryExists(userID int64, url string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM urls WHERE user_id = $1 AND url = $2)`,
		userID, url,
	).Scan(&exists)
	return exists, err
}

// CreateUrlEntry inserts a new URL entry linked to a user.
//
// Parameters:
//   - entry: Url model containing the user ID, the URL, and optionally the original file line.
//
// Returns:
//   - ID of the new URL entry.
//   - error if the insert fails.
func (r *CombolistRepository) CreateUrlEntry(entry models.Urls) (int64, error) {
	var id int64
	query := `
		INSERT INTO urls (user_id, url, file_line)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRow(query, entry.UserID, entry.URL, 0).Scan(&id)
	return id, err
}

// UrlExistsByCredentialsAndFile checks if a URL already exists for the same username + password in the same file.
func (r *CombolistRepository) UrlExistsByCredentialsAndFile(fileID int64, username, password, url string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM urls u
			INNER JOIN users usr ON u.user_id = usr.id
			WHERE usr.file_id = $1 AND usr.username = $2 AND usr.password = $3 AND u.url = $4
		)
	`
	err := r.db.QueryRow(query, fileID, username, password, url).Scan(&exists)
	return exists, err
}
