package utils

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// ValidatePassword validates password strength
// Requirements: 8+ chars, uppercase, lowercase, number
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}

// ValidateEmail validates email format and length
func ValidateEmail(email string) error {
	if len(email) > 255 {
		return errors.New("email must be 255 characters or less")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateUsername validates username format and length
// Requirements: 3-50 chars, alphanumeric + underscore/hyphen
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if len(username) > 50 {
		return errors.New("username must be 50 characters or less")
	}

	// Only alphanumeric, underscore, and hyphen allowed
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if err != nil {
		return errors.New("username validation failed")
	}
	if !matched {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// SanitizeHTML removes HTML tags and dangerous characters
func SanitizeHTML(input string) string {
	// Remove HTML tags
	re := regexp.MustCompile("<[^>]*>")
	sanitized := re.ReplaceAllString(input, "")

	// Remove script-related content
	re = regexp.MustCompile("(?i)<script[^>]*>.*?</script>")
	sanitized = re.ReplaceAllString(sanitized, "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// ValidateRecipeTitle validates recipe title
func ValidateRecipeTitle(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return errors.New("recipe title is required")
	}

	if len(title) > 200 {
		return errors.New("recipe title must be 200 characters or less")
	}

	return nil
}

// ValidateRecipeDescription validates recipe description
func ValidateRecipeDescription(description string) error {
	if len(description) > 2000 {
		return errors.New("recipe description must be 2000 characters or less")
	}

	return nil
}

// ValidateIngredients validates recipe ingredients list
func ValidateIngredients(ingredients []string) error {
	if len(ingredients) > 50 {
		return errors.New("recipe can have at most 50 ingredients")
	}

	for i, ingredient := range ingredients {
		if len(ingredient) > 200 {
			return errors.New("ingredient " + string(rune(i+1)) + " must be 200 characters or less")
		}
	}

	return nil
}

// ValidateInstructions validates recipe instructions list
func ValidateInstructions(instructions []string) error {
	if len(instructions) > 100 {
		return errors.New("recipe can have at most 100 instruction steps")
	}

	for i, instruction := range instructions {
		if len(instruction) > 1000 {
			return errors.New("instruction step " + string(rune(i+1)) + " must be 1000 characters or less")
		}
	}

	return nil
}

// ValidateMealType validates meal type enum
func ValidateMealType(mealType string) error {
	validTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
	}

	if !validTypes[strings.ToLower(mealType)] {
		return errors.New("meal type must be one of: breakfast, lunch, dinner, snack")
	}

	return nil
}

// ValidateDateFormat validates YYYY-MM-DD date format
func ValidateDateFormat(date string) error {
	matched, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}$", date)
	if err != nil {
		return errors.New("date validation failed")
	}
	if !matched {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	return nil
}
