package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/DataDog/go-sqllexer"
	"strings"
)

func NormalizeQuery(query string) string {
	query = strings.TrimSpace(query)
	obfuscator := sqllexer.NewObfuscator()
	obfuscated := obfuscator.Obfuscate(query)

	return obfuscated
}

func NormalizeAndHashQuery(query string) (string, string) {
	normalized := NormalizeQuery(query)

	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(normalized))
	hash := sha256Hash.Sum(nil)

	return normalized, hex.EncodeToString(hash)
}
