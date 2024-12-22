package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

// GenerateSha1Hash generates the SHA-1 hash of the given text.
//
// Parameters:
// - text: A string containing the text to be hashed.
//
// Returns:
// - hash: A string containing the SHA-1 hash of the input text.
//
// The function creates a new SHA-1 hash instance, writes the provided text as a byte slice,
// and then encodes the resulting hash into a hexadecimal string. The hash is returned as a string.
func GenerateSha1Hash(text string) (hash string) {
    sha1 := sha1.New()
    sha1.Write([]byte(text))
    hash = hex.EncodeToString(sha1.Sum(nil))

    return
}
