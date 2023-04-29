package slugify

import (
	"regexp"
	"strings"
)

// slugify converts a string to a slug.
func Slugify(input string) string {
	// Remove any non-alphanumeric or non-hyphen characters and replace with hyphens
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")
	if err != nil {
		panic(err)
	}

	// Replace any multiple consecutive hyphens with a single hyphen
	input = reg.ReplaceAllString(strings.TrimSpace(input), "-")

	// Convert the string to lowercase
	return strings.ToLower(input)
}
