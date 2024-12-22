package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sentielxx/searchx-indexer/internal/config"
	"github.com/sentielxx/searchx-indexer/internal/db"
	"github.com/sentielxx/searchx-indexer/internal/runner"
	"github.com/sentielxx/searchx-indexer/pkg/hash"
	"github.com/sentielxx/searchx-indexer/pkg/utils"
	"github.com/sentielxx/searchx-indexer/process"
)

func init() {
	runner.Banner()
}

func main() {
	configx, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("[!] the file 'config.yaml' was not found.\n")
	}

	filedir := flag.String("file", "", "path to the log file.")
	description := flag.String("description", "not information", "optional description for the logins.")
	source := flag.String("source", "not information", "source dataleak")

	flag.Parse()

	if *filedir == "" {
		log.Fatal("[!] you did not provide the path to the log file, use the '-file' flag.")
	}

	fileInfo, err := os.Stat(*filedir)
	if os.IsNotExist(err) {
		log.Fatalf("[!] the log file '%s' was not found.%e\n", *filedir, err)
	}

	database := db.NewDatabase(configx.Database.Port, configx.Database.Host, configx.Database.Username, configx.Database.Password, configx.Database.Dbname, configx.Database.Sslmode)

	cursor, err := database.Connect()
	if err != nil {
		log.Fatalf("[!] %s\n", err.Error())
	}
	defer database.Close()

	file := db.ModelComboListFile{
		Name:                  fileInfo.Name(),
		Hash:                  hash.GenerateSha1Hash(fileInfo.Name() + fmt.Sprint(fileInfo.Size())),
		Size:                  fileInfo.Size(),
		Type:                  strings.ReplaceAll(filepath.Ext(*filedir), ".", ""),
		Description:           *description,
		Source:                *source,
		CreatedAt:             fileInfo.ModTime().Format("2006-01-02"),
		ProcessedEntriesCount: 0,
		ProcessedAt:           time.Now().Format("2006-01-02"),
	}

	log.Println("file hash       : ", file.Hash)
	log.Println("file name       : ", file.Name)
	log.Println("file description: ", file.Description)
	log.Println("file source     : ", file.Source)
	log.Println("file type       : ", file.Type)
	log.Println("file size       : ", utils.ConvertFileSize(file.Size))
	log.Println("file created at : ", file.CreatedAt)
	println()

	err = process.ComboListProcessDataleak(cursor, &file)

	if err != nil {
		log.Fatalln(err)
	}
}
