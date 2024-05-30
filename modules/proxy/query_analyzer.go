package proxy

import (
	"strings"
	"unicode"
)

// CommandType represents the type of SQL command.
type CommandType int

const (
	Unknown     CommandType = iota // Unknown command type.
	SetNames                       // SET NAMES command type.
	UseDatabase                    // USE DATABASE command type.
)

// SQLCommand represents a parsed SQL command.
type SQLCommand struct {
	Type  CommandType // Type of the SQL command.
	Value string      // Value associated with the SQL command.
}

// Analyze analyzes an SQL query string and returns an SQLCommand with the detected type and value.
func Analyze(query string) SQLCommand {
	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return SQLCommand{Type: Unknown}
	}

	if command, found := AnalyzeSetNamesCharset(query); found {
		return command
	}

	if command, found := AnalyzeUseDatabase(query); found {
		return command
	}

	return SQLCommand{Type: Unknown}
}

// AnalyzeSetNamesCharset extracts the charset from a SET NAMES SQL query.
// Returns an SQLCommand with the charset if the query matches the SET NAMES pattern.
func AnalyzeSetNamesCharset(query string) (SQLCommand, bool) {
	const setNamesPrefix = "SET NAMES"
	queryLen := len(query)
	prefixLen := len(setNamesPrefix)

	i, j := 0, 0

	// Skip leading spaces
	for i < queryLen && unicode.IsSpace(rune(query[i])) {
		i++
	}

	// Compare SET NAMES case-insensitively
	for j < prefixLen && i < queryLen {
		if unicode.ToUpper(rune(query[i])) != rune(setNamesPrefix[j]) {
			return SQLCommand{}, false
		}
		i++
		j++
	}

	// Ensure the entire prefix was matched
	if j != prefixLen {
		return SQLCommand{}, false
	}

	// Skip any spaces after SET NAMES
	for i < queryLen && unicode.IsSpace(rune(query[i])) {
		i++
	}

	// The remaining part is the charset value
	if i < queryLen {
		return SQLCommand{Type: SetNames, Value: query[i:]}, true
	}

	return SQLCommand{}, false
}

// AnalyzeUseDatabase extracts the database name from a USE SQL query.
// Returns an SQLCommand with the database name if the query matches the USE pattern.
func AnalyzeUseDatabase(query string) (SQLCommand, bool) {
	const usePrefix = "USE"
	queryLen := len(query)
	prefixLen := len(usePrefix)

	i, j := 0, 0

	// Skip leading spaces
	for i < queryLen && unicode.IsSpace(rune(query[i])) {
		i++
	}

	// Compare USE case-insensitively
	for j < prefixLen && i < queryLen {
		if unicode.ToUpper(rune(query[i])) != rune(usePrefix[j]) {
			return SQLCommand{}, false
		}
		i++
		j++
	}

	// Ensure the entire prefix was matched
	if j != prefixLen {
		return SQLCommand{}, false
	}

	// Skip any spaces after USE
	for i < queryLen && unicode.IsSpace(rune(query[i])) {
		i++
	}

	// The remaining part is the database name
	if i < queryLen {
		return SQLCommand{Type: UseDatabase, Value: query[i:]}, true
	}

	return SQLCommand{}, false
}
