package handlers

import (
	"encoding/json"
	"log"
	"strings"

	models "searchx-indexer/models/combolist"
	repository "searchx-indexer/repository/combolist"
)

// HandleCombolistMessage receives and processes messages related to combolist operations.
// It expects a message with a "type" field that determines the type of operation.
// Currently supports:
//   - "file"   → Processes file metadata
//   - "entrie" → Processes a username/password/url entry tied to a file
//
// Example expected input structure:
//
//	map[string]interface{}{
//	  "type": "file" or "entrie",
//	  ...
//	}
func HandleCombolistMessage(data map[string]interface{}, db *repository.CombolistRepository) {
	typex, ok := data["type"].(string)
	if !ok {
		log.Println("[!] Invalid or missing type in incoming message")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("[!] Failed to marshal data to JSON:", err)
		return
	}

	switch typex {
	case "file":
		handleCombolistFile(jsonData, db)

	case "entrie":
		handleCombolistEntrie(jsonData, db)

	default:
		log.Printf("[!] Unknown type received: %s", typex)
	}
}

// handleCombolistFile processes file metadata messages and saves them if they don't exist.
// It checks for the presence of a valid hash, name, and size before persisting.
// The file is uniquely identified by its hash.
func handleCombolistFile(jsonData []byte, db *repository.CombolistRepository) {
	var file models.SocketFile
	if err := json.Unmarshal(jsonData, &file); err != nil {
		log.Println("[!] Failed to parse file structure:", err)
		return
	}

	model := file.Raw.Files
	hash := model.Hash

	if hash == "" || model.Name == "" || model.Size == 0 {
		log.Println("[!] Missing required file fields (hash, name or size)")
		return
	}

	model.AgentKey = file.AgentKey

	// Check if this file already exists in the DB based on its hash
	exists, err := db.FileExists(hash)
	if err != nil {
		log.Println("[!] Error checking file existence:", err)
		return
	}
	if exists {
		log.Println("[-] File already exists in the database (hash check)")
		return
	}

	// Insert file into the database
	if _, err := db.CreateFile(model); err != nil {
		log.Println("[!] Error creating combolist file:", err)
	}
}

// handleCombolistEntrie processes a single username/password/URL entry.
// It ensures the user exists only for the specific file, and avoids duplicate URL insertions
// for the same user *within that file*.
func handleCombolistEntrie(jsonData []byte, db *repository.CombolistRepository) {
	var entry models.SocketEntrie
	if err := json.Unmarshal(jsonData, &entry); err != nil {
		log.Println("[!] Failed to parse entry structure:", err)
		return
	}

	// === 1. Extract and validate required fields
	username := strings.TrimSpace(entry.Raw.Users.Username)
	password := strings.TrimSpace(entry.Raw.Users.Password)
	url := strings.TrimSpace(entry.Raw.Urls.URL)
	hash := strings.TrimSpace(entry.Hash)

	if username == "" || password == "" {
		log.Println("[!] Missing username or password in entry")
		return
	}
	if url == "" {
		log.Println("[!] Missing URL in entry")
		return
	}
	if hash == "" {
		log.Println("[!] Missing file hash in entry")
		return
	}

	// === 2. Lookup file based on hash
	fileID, err := db.GetFileByHash(hash)
	if err != nil || fileID == 0 {
		log.Printf("[!] File not found (hash: %s) — entry processing aborted", hash)
		return
	}

	// === 3. Check if URL already exists for the given (file + username + password)
	// This ensures the same user credentials *in the same file* don’t get duplicate URLs
	exists, err := db.UrlExistsByCredentialsAndFile(fileID, username, password, url)
	if err != nil {
		log.Printf("[!] Error checking URL entry with context (fileID=%d, username=%s): %v", fileID, username, err)
		return
	}
	if exists {
		log.Printf("[-] URL already exists for username [%s] in file ID %d", username, fileID)
		return
	}

	// === 4. Check if user already exists in this file
	// NOTE: Critical to match not only by username/password but also by file ID
	userID, err := db.GetUserIDByCredentialsInFile(fileID, username, password)
	if err != nil {
		log.Printf("[!] Error checking user entry (username: %s): %v", username, err)
		return
	}

	if userID > 0 {
		// === 5. User exists: associate new URL to existing user
		if _, err := db.CreateUrlEntry(models.Urls{
			UserID: userID,
			URL:    url,
		}); err != nil {
			log.Printf("[!] Error inserting URL entry for existing user ID %d: %v", userID, err)
		}
		return
	}

	// === 6. User doesn't exist: create user
	newUserID, err := db.CreateUserEntry(models.Users{
		FileID:   fileID,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Printf("[!] Error creating user entry (username: %s): %v", username, err)
		return
	}

	// === 7. Associate URL to the newly created user
	if _, err := db.CreateUrlEntry(models.Urls{
		UserID: newUserID,
		URL:    url,
	}); err != nil {
		log.Printf("[!] Error inserting URL entry for new user ID %d: %v", newUserID, err)
	}
}
