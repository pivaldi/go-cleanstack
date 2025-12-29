package migrations

import (
	"fmt"
	"strings"
)

// ValidateDescription validates a migration description.
// Returns error if description is invalid.
func ValidateDescription(description string) error {
	if description == "" {
		return fmt.Errorf("description is required")
	}

	if len(description) < 3 || len(description) > 100 {
		return fmt.Errorf("description must be between 3 and 100 characters")
	}

	return nil
}

// ToCamelCase converts a string to CamelCase for Go function names.
// Handles dashes, underscores, and spaces as separators.
func ToCamelCase(s string) string {
	// Replace dashes and underscores with spaces for uniform processing
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Split by spaces
	words := strings.Fields(s)

	// Capitalize first letter of each word
	for i, word := range words {
		if word != "" {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}
