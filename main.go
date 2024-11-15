package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sentielxx/xlog-process/database"
	"github.com/sentielxx/xlog-process/process"
	"github.com/sentielxx/xlog-process/runner"
)

func init() {
	runner.Banner()
}

func main() {
	file := flag.String("file", "", "path to the log file.")
	databaseFile := flag.String("database", "database/xlog.db", "path to the sqlite database.")
	description := flag.String("description", "not information", "optional description for the logins.")
	disableSecurity := flag.Bool("disable-security", false, "disable options for security.")
	enableMemory := flag.Bool("enable-temp-memory", false, "")
	blockSize := flag.Int64("block-size", 1000, "set block size to process.")

	flag.Parse()

	if *file == "" {
		print("[!] you did not provide the path to the log file, use the '-file' flag.")
	}

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		fmt.Printf("[!] the log file '%s' was not found.\n", *file)
		os.Exit(1)
	}

	conn, err := database.Connect(*databaseFile)
	if err != nil {
		fmt.Printf("[!] the log file '%s' was not found.\n", *file)
		os.Exit(1)
	}
	defer conn.Close()

	if *disableSecurity {
		if err = database.DisableSecuritySettings(conn); err != nil {
			print("[!] error disable security:", err, "\n")
			os.Exit(1)
		}
	}

	if *enableMemory {
		if err = database.SetTempStoreMemorySpeed(conn); err != nil {
			print("[!] error set temp memory:", err, "\n")
			os.Exit(1)
		}
	}

	if err = database.CredentialsCreateTable(conn); err != nil {
		print("[!] error create table:", err, "\n")
	}

	linestotal, linesprocessed := process.ProcessData(conn, *file, *description, *blockSize)

	println("[+] total lines:", linestotal)
	println("[+] total lines processed:", linesprocessed)
	println("[+] processing completed successfully!")
}
