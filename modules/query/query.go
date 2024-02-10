package query

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
	obfuscated := NormalizeQuery(query)

	hasher := sha256.New()
	hasher.Write([]byte(obfuscated))
	hash := hasher.Sum(nil)

	return obfuscated, hex.EncodeToString(hash)
}
