package db

func defaultConfig(dataconfig *Database) {
	if dataconfig.Host == "" {
		dataconfig.Host = "127.0.0.1"
	}
	if dataconfig.Port == "" {
		dataconfig.Port = "5432"
	}
	if dataconfig.Username == "" {
		dataconfig.Username = "postgres"
	}
	if dataconfig.Password == "" {
		dataconfig.Password = ""
	}
	if dataconfig.SSL == "" {
		dataconfig.SSL = "disable"
	}
}
