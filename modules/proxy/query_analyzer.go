package proxy

import (
	"unicode"
)

// CommandType represents the type of SQL command.
type CommandType int

const (
	Unknown     CommandType = iota // Unknown command type.
	SetNames                       // SET NAMES command type.
	UseDatabase                    // USE DATABASE command type.
)

const (
	SetNamesPrefix = "SET NAMES"
	UsePrefix      = "USE"
)

// SQLCommand represents a parsed SQL command.
type SQLCommand struct {
	Type  CommandType // Type of the SQL command.
	Value string      // Value associated with the SQL command.
}

// Analyze analyzes an SQL query string and returns an SQLCommand with the detected type and value.
func Analyze(query string) SQLCommand {
	if len(query) == 0 {
		return SQLCommand{Type: Unknown}
	}

	if command, found := analyzePrefixedCommand(query, SetNamesPrefix, SetNames); found {
		return command
	}

	if command, found := analyzePrefixedCommand(query, UsePrefix, UseDatabase); found {
		return command
	}

	return SQLCommand{Type: Unknown}
}

// analyzePrefixedCommand extracts the value from an SQL query if it matches the given prefix case-insensitively.
// Returns an SQLCommand with the given command type and the value if the query matches the pattern.
func analyzePrefixedCommand(query, prefix string, commandType CommandType) (SQLCommand, bool) {
	queryLen := len(query)
	prefixLen := len(prefix)

	i, j := 0, 0

	// Skip leading spaces
	for i < queryLen && unicode.IsSpace(rune(query[i])) {
		i++
	}

	// Compare prefix case-insensitively
	for j < prefixLen && i < queryLen {
		if unicode.ToUpper(rune(query[i])) != rune(prefix[j]) {
			return SQLCommand{}, false
		}
		i++
		j++
	}

	// Ensure the entire prefix was matched
	if j != prefixLen {
		return SQLCommand{}, false
	}

	// Skip any spaces after the prefix
	for i < queryLen && unicode.IsSpace(rune(query[i])) {
		i++
	}

	// The remaining part is the value
	if i < queryLen {
		return SQLCommand{Type: commandType, Value: query[i:]}, true
	}

	return SQLCommand{}, false
}
