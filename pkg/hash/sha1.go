package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func TextToSha1(text string) (hash string) {
	sha1 := sha1.New()
	sha1.Write([]byte(text))
	hash = hex.EncodeToString(sha1.Sum(nil))

	return
}
