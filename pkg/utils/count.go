package utils

import (
	"bufio"
	"fmt"
	"os"
)

// CountLinesInFile counts the total number of lines in a given file.
//
// Parameters:
// - filePath: A string representing the path to the file whose lines are to be counted.
//
// Returns:
// - lineCount: The total number of lines in the file as an int64.
// - err: An error object that will be non-nil if there was an issue opening the file or reading its contents.
//
// The function opens the file and reads it line by line using a scanner, 
// incrementing the line counter for each line it encounters. It checks for errors 
// during scanning and handles them accordingly. If no errors occur, it returns the 
// total line count and a nil error.
func CountLinesInFile(filePath string) (lineCount int64, err error) {
    fileHandle, err := os.Open(filePath)
    if err != nil {
        return 0, fmt.Errorf("error opening file '%s': %v", filePath, err)
    }
    defer fileHandle.Close()

    lineScanner := bufio.NewScanner(fileHandle)

    for lineScanner.Scan() {
        lineCount++
    }

    if err := lineScanner.Err(); err != nil {
        return 0, fmt.Errorf("error reading file '%s': %v", filePath, err)
    }

    return lineCount, nil
}
