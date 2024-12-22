package utils

import "fmt"

// ConvertFileSize converts a file size in bytes to a human-readable unit,
// such as Bytes, KB, MB, or GB. It returns a formatted string with the
// converted size and the appropriate unit.
//
// Example:
//  ConvertFileSize(1234567890) // "1.15 GB"
//
// Parameters:
//  sizeInBytes: The file size in bytes (int64).
//
// Returns:
//  A string representing the file size, formatted with two decimal places
//  and the corresponding unit (Bytes, KB, MB, or GB).
func ConvertFileSize(sizeInBytes int64) string {
	var size float64
	var unit string

	if sizeInBytes < 1024 {
		size = float64(sizeInBytes)
		unit = "Bytes"
	} else if sizeInBytes < 1024*1024 {
		size = float64(sizeInBytes) / 1024
		unit = "KB"
	} else if sizeInBytes < 1024*1024*1024 {
		size = float64(sizeInBytes) / (1024 * 1024)
		unit = "MB"
	} else {
		size = float64(sizeInBytes) / (1024 * 1024 * 1024)
		unit = "GB"
	}

	return fmt.Sprintf("%.2f %s", size, unit)
}
