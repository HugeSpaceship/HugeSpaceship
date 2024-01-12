package validation

import "regexp"

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func IsHashValid(hash string) bool {
	return !nonAlphanumericRegex.MatchString(hash) && len(hash) > 40
}
