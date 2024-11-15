package process

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/sentielxx/xlog-process/database"
)

func processBlock(conn *sql.DB, mtx *sync.Mutex, linesprocessed *int64, blocklines []string, filename string, filesize int64, description string) {
	regex := `([a-zA-Z0-9+.-]+://[^:/\s]+(?:/[^:\s]*)?):([^:]+):([^\n]*)`

	re := regexp.MustCompile(regex)

	for _, line := range blocklines {
		mtx.Lock()

		matches := re.FindStringSubmatch(strings.TrimSpace(line))

		if matches != nil {
			err := database.CredentialsInsertTable(conn, matches[1], matches[2], matches[3], description, filename, filesize)
			if err != nil {
				println("[!] error in insert table credentials:", err)
			}

			*linesprocessed++
		}

		mtx.Unlock()
	}
}

func ProcessData(conn *sql.DB, filename string, description string, blocksize int64) (int64, int64) {
	fileinfo, err := os.Stat(filename)

	if os.IsNotExist(err) {
		fmt.Printf("[!] the file '%s' was not found.\n", filename)
		os.Exit(1)
	}
	filesize := fileinfo.Size()

	filebytes, err := os.Open(filename)
	if err != nil {
		fmt.Printf("[!] the file '%s' was not found.\n", filename)
		os.Exit(1)
	}
	defer filebytes.Close()

	var linestotal int64
	var linesprocessed int64

	scanner := bufio.NewScanner(filebytes)
	for scanner.Scan() {
		linestotal++
	}
	filebytes.Seek(0, 0)

	var blocklines []string

	var wg sync.WaitGroup
	var mtx sync.Mutex

	threadscontrol := make(chan struct{}, 2)

	scanner = bufio.NewScanner(filebytes)
	for scanner.Scan() {
		filelinestring := scanner.Text()

		blocklines = append(blocklines, filelinestring)

		if int64(len(blocklines)) >= blocksize {
			threadscontrol <- struct{}{}
			wg.Add(1)
			go func(blocklines []string) {
				defer wg.Done()
				processBlock(conn, &mtx, &linesprocessed, blocklines, filename, filesize, description)

				mtx.Lock()
				linesprocessed += int64(len(blocklines))
				mtx.Unlock()

				progress := float64(linesprocessed) / float64(linestotal) * 100
				fmt.Printf("\r[%.2f%%] lines total: %d | lines processed: %d | lines ignored: %d", progress, linestotal, linesprocessed, linestotal-linesprocessed)
				<-threadscontrol
			}(blocklines)

			blocklines = []string{}
		}
	}

	if len(blocklines) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processBlock(conn, &mtx, &linesprocessed, blocklines, filename, filesize, description)
			mtx.Lock()
			linesprocessed += int64(len(blocklines))
			mtx.Unlock()
			progress := float64(linesprocessed) / float64(linestotal) * 100
			fmt.Printf("\r[%.2f%%] lines total: %d | lines processed: %d | lines ignored: %d", progress, linestotal, linesprocessed, linestotal-linesprocessed)
		}()
	}

	wg.Wait()

	return linestotal, linesprocessed
}
