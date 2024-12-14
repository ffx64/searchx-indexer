package process

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sentielxx/xlog-process/internal/db"
	"github.com/sentielxx/xlog-process/pkg/hash"
)

func StealerLogsProcessData(cursor *db.Database, filename, description, source string) error {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("o arquivo '%s' não foi encontrado", filename)
	}

	fileBytes, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo '%s': %v", filename, err)
	}
	defer fileBytes.Close()

	fileHash := hash.TextToSha1(fileInfo.Name())
	fileLogDatabaseExists, err := cursor.FileLogExists(fileHash, fileInfo.Name())
	if err != nil {
		return fmt.Errorf("erro ao verificar se o log existe: %v", err)
	}
	if fileLogDatabaseExists {
		log.Printf("[!] Log de arquivo '%s' já existe no banco de dados", filename)
		return nil
	}

	fileModel := db.StealerFileLogModel{
		FileName:        fileInfo.Name(),
		FileSize:        fileInfo.Size(),
		FileType:        fileInfo.Mode().Type().String(),
		FileDescription: description,
		FileHash:        fileHash,
		Source:          source,
		CreatedAt:       fileInfo.ModTime().Format("2006-01-02"),
		ProcessedAt:     time.Now().Format("2006-01-02"),
	}

	fileModel.FileID, err = cursor.InsertStealerFile(fileModel)
	if err != nil {
		return fmt.Errorf("erro ao inserir log de arquivo '%s': %v", filename, err)
	}

	if err := processFileLines(fileBytes, cursor, fileModel); err != nil {
		return err
	}

	log.Printf("[*] Processamento do arquivo '%s' concluído com sucesso", filename)
	return nil
}

func processFileLines(fileBytes *os.File, cursor *db.Database, fileModel db.StealerFileLogModel) error {
	regx := regexp.MustCompile(`([a-zA-Z0-9+.-]+://[^:/\s]+(?:/[^:\s]*)?):([^:]+):([^\n]*)`)

	var linesTotal, linesProcessed int64

	for bufio.NewScanner(fileBytes).Scan() {
		linesTotal++
	}
	fileBytes.Seek(0, 0)

	scanner := bufio.NewScanner(fileBytes)
	for scanner.Scan() {
		fileLine := strings.TrimSpace(scanner.Text())

		if match := regx.FindStringSubmatch(fileLine); match != nil {
			entryModel := db.StealerEntrieLogModel{
				URL:         match[1],
				Username:    match[2],
				Password:    match[3],
				CreatedAt:   fileModel.CreatedAt,
				ProcessedAt: time.Now().Format("2006-01-02"),
			}

			if err := cursor.InsertStealerEntrie(entryModel, fileModel.FileID); err != nil {
				log.Printf("[!] Erro ao inserir entrada de stealer: %v", err)
			}

			linesProcessed++
		}

		if linesTotal > 0 {
			progress := float64(linesProcessed) / float64(linesTotal) * 100
			fmt.Printf("\r[%.2f%%] Total de linhas: %d | Processadas: %d | Ignoradas: %d", progress, linesTotal, linesProcessed, linesTotal-linesProcessed)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erro ao ler as linhas do arquivo: %v", err)
	}

	log.Printf("[*] %d linhas processadas", linesProcessed)
	return nil
}
