package proxy

import (
	"go-proxy/modules/log"
	"go.uber.org/zap"
	"strings"
	"unicode"
)

// transactionCommands defines the mapping of SQL transaction commands to their respective handlers.
var (
	transactionCommands = map[string]CommandHandler{
		"START TRANSACTION": handleBeginTransaction,
		"BEGIN":             handleBeginTransaction,
		"COMMIT":            handleCommit,
		"ROLLBACK":          handleRollback,
	}
)

// CommandHandler defines a function type for handling SQL commands.
type CommandHandler func(proxy *ProxyHandler, query string)

// AnalyzeTransaction analyzes a SQL query to determine if it is a transaction command
// and invokes the corresponding handler if a match is found.
func AnalyzeTransaction(proxy *ProxyHandler, query string) {
	query = strings.TrimSpace(query)

	for cmd, handler := range transactionCommands {
		if caseInsensitiveHasPrefix(query, cmd) {
			handler(proxy, query)
			break
		}
	}
}

// handleBeginTransaction handles the start of a transaction.
func handleBeginTransaction(proxy *ProxyHandler, query string) {
	log.Logger.Debug("Begin Transaction", zap.String("query", query))
	proxy.transaction = true
	proxy.sendInTransaction = true
}

// handleCommit handles the commit of a transaction.
func handleCommit(proxy *ProxyHandler, query string) {
	log.Logger.Debug("Commit Transaction", zap.String("query", query))
	proxy.transaction = false
	proxy.sendInTransaction = true
}

// handleRollback handles the rollback of a transaction.
func handleRollback(proxy *ProxyHandler, query string) {
	log.Logger.Debug("Rollback Transaction", zap.String("query", query))
	proxy.transaction = false
	proxy.sendInTransaction = true
}

// caseInsensitiveHasPrefix checks if the query starts with the specified prefix, ignoring case.
func caseInsensitiveHasPrefix(query, prefix string) bool {
	i, j := 0, 0
	queryLen, prefixLen := len(query), len(prefix)

	for i < queryLen && j < prefixLen {
		// Skip spaces in query
		for i < queryLen && unicode.IsSpace(rune(query[i])) {
			i++
		}
		// Skip spaces in prefix
		for j < prefixLen && unicode.IsSpace(rune(prefix[j])) {
			j++
		}
		if i < queryLen && j < prefixLen {
			if unicode.ToUpper(rune(query[i])) != unicode.ToUpper(rune(prefix[j])) {
				return false
			}
			i++
			j++
		}
	}

	// Skip trailing spaces in prefix
	for j < prefixLen && unicode.IsSpace(rune(prefix[j])) {
		j++
	}

	return j == prefixLen
}
