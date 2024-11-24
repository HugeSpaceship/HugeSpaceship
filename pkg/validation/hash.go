package validation

import "regexp"

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// IsHashValid validates that a resource doesn't contain invalid characters and isn't too long
// It could possibly be smarter
func IsHashValid(hash string) (bool, string) {
	if nonAlphanumericRegex.MatchString(hash) || len(hash) > 40 {
		return false, ""
	} else {
		return true, hash
	}

}
