package sqlsafety

import (
	"fmt"
	"regexp"
	"strings"
)

// SQLSafetyChecker provides utilities for SQL injection prevention
type SQLSafetyChecker struct {
	// Valid column names for ordering
	validColumns map[string]bool
	// Regex for validating input
	identifierRegex *regexp.Regexp
	numberRegex     *regexp.Regexp
}

// NewSQLSafetyChecker creates a new instance with the given valid columns
func NewSQLSafetyChecker(validColumns []string) *SQLSafetyChecker {
	columnsMap := make(map[string]bool)
	for _, col := range validColumns {
		columnsMap[col] = true
	}

	return &SQLSafetyChecker{
		validColumns:    columnsMap,
		identifierRegex: regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
		numberRegex:     regexp.MustCompile(`^[0-9]+$`),
	}
}

// ValidateOrderBy checks if the order by column is valid
func (s *SQLSafetyChecker) ValidateOrderBy(column string) error {
	if !s.validColumns[column] {
		return fmt.Errorf("invalid order by column: %s", column)
	}
	return nil
}

// ValidateDirection checks if the sort direction is valid
func (s *SQLSafetyChecker) ValidateDirection(direction string) error {
	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		return fmt.Errorf("invalid sort direction: %s", direction)
	}
	return nil
}

// ValidateIdentifier checks if the identifier is safe
func (s *SQLSafetyChecker) ValidateIdentifier(identifier string) error {
	if !s.identifierRegex.MatchString(identifier) {
		return fmt.Errorf("invalid identifier format: %s", identifier)
	}
	return nil
}

// ValidateNumber checks if the string is a valid number
func (s *SQLSafetyChecker) ValidateNumber(number string) error {
	if number != "" && !s.numberRegex.MatchString(number) {
		return fmt.Errorf("invalid number format: %s", number)
	}
	return nil
}

// SanitizeSearchTerm removes potentially dangerous characters from search terms
func (s *SQLSafetyChecker) SanitizeSearchTerm(term string) string {
	// Remove any SQL special characters
	dangerous := []string{"'", "\"", ";", "--", "/*", "*/", "xp_"}
	safe := term
	for _, d := range dangerous {
		safe = strings.ReplaceAll(safe, d, "")
	}
	return safe
}
