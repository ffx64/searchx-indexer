package combolist

import models "searchx-indexer/models/combolist"

// CombolistRepositoryInterface defines all necessary methods to interact with the combolist data
type CombolistRepositoryInterface interface {
	// File-related methods
	FileExists(hash string) (bool, error)
	GetFileByHash(hash string) (int64, error)
	CreateFile(file models.Files) (int64, error)

	// User-related methods
	GetUserIDByCredentialsInFile(fileID int64, username, password string) (int64, error)
	CreateUserEntry(user models.Users) (int64, error)

	// URL-related methods
	UrlEntryExists(userID int64, url string) (bool, error)
	CreateUrlEntry(entry models.Urls) (int64, error)
}
