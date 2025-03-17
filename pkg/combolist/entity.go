package combolist

// EntityCombolistFile represents a file combolist in the database, containing metadata about a file.
type EntityCombolistFile struct {
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

// EntityCombolistUserEntrie represents an entry in the combolist file, containing sensitive data.
type EntityCombolistUserEntrie struct {
	ID        int64
	FileID    int64
	Username  string // Username is the username for the entry.
	Password  string // Password is the password associated with the entry.
	CreatedAt string // CreatedAt is the timestamp when the entry was created.
}

type EntityCombolistUrlEntrie struct {
	ID        int64
	UserID    int64
	URL       string // URL is the URL associated with the combolist entry.
	FileLine  int64  // FileLine is the line number in the combolist file where this entry was found.
	CreatedAt string
}
