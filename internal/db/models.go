package db

import "database/sql"

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSL      string

	cursor *sql.DB
}

type StealerFileLogModel struct {
	FileID                int64
	FileName              string
	FileSize              int64
	FileHash              string
	Status                string
	Source                string
	FileType              string
	FileDescription       string
	ProcessedEntriesCount int64
	CreatedAt             string
	ProcessedAt           string
}

type StealerEntrieLogModel struct {
	URL         string
	Username    string
	Password    string
	CreatedAt   string
	ProcessedAt string
}
