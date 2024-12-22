package process

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sentielxx/searchx-indexer/internal/db"
	"github.com/sentielxx/searchx-indexer/pkg/utils"
)

// ComboListProcessDataleak processes a ComboList data leak file, extracting
// entries (URL, Username, Password) and storing them in the database.
// It also checks for duplicate entries to avoid redundant inserts.
//
// This function expects the ComboList file to be in a specific format with
// URL:Username:Password entries. It reads through the file line by line,
// and for each line that matches the expected pattern, it inserts the
// entry into the database if it does not already exist.
//
// Parameters:
//   - database: A reference to the database object used to interact with
//     the database and perform CRUD operations.
//   - file: The ComboList file containing the data leak entries.
//
// Returns:
//   - An error if something goes wrong during processing or database operations.
func ComboListProcessDataleak(database *db.Database, file *db.ModelComboListFile) error {
	regx := regexp.MustCompile(`([a-zA-Z0-9+.-]+://[^:/\s]+(?:/[^:\s]*)?):([^:]+):([^\n]*)`)

	exists, err := database.DoesComboListFile(file.Hash)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("file with hash %s has already been processed and is stored in the database", file.Hash)
	}

	file.ID, err = database.CreateComboListFile(*file)
	if err != nil {
		return err
	}

	linesTotal, err := utils.CountLinesInFile(file.Name)
	if err != nil {
		return err
	}

	fileBytes, err := os.Open(file.Name)
	if err != nil {
		return fmt.Errorf("error opening file '%s': %v", file.Name, err)
	}
	defer fileBytes.Close()

	var linesProcessed int64

	scanner := bufio.NewScanner(fileBytes)
	for scanner.Scan() {
		match := regx.FindStringSubmatch(strings.TrimSpace(scanner.Text()))
		if match == nil {
			continue // Skip lines that don't match
		}

		entrie := db.ModelComboListEntrie{
			URL:         match[1],
			Username:    match[2],
			Password:    match[3],
			FileLine:    linesProcessed,
			CreatedAt:   file.CreatedAt,
			ProcessedAt: time.Now().Format("2006-01-02"),
		}

		exists, err := database.DoesComboListEntrie(entrie.URL, entrie.Username, entrie.Password)
		if err != nil {
			return err
		}

		if exists { // Skip inserting if the entry already exists
			log.Printf("[*] Entry already exists for URL: %s, Username: %s, skipping insert.", entrie.URL, entrie.Username)
			continue
		}

		_, err = database.CreateComboListEntrie(entrie, file.ID)
		if err != nil {
			return err
		}

		linesProcessed++

		fmt.Printf("\r[%.2f%%] | Total lines: %d | Processed: %d | Ignored: %d",
			float64(linesProcessed)/float64(linesTotal)*100,
			linesTotal,
			linesProcessed,
			linesTotal-linesProcessed)
	}

	println("")
	log.Println("[+] all records have been successfully added to the database.")

	return nil
}
