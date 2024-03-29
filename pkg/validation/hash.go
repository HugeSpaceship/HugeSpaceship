package validation

import "regexp"

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func IsHashValid(hash string) (bool, string) {
	if nonAlphanumericRegex.MatchString(hash) || len(hash) > 40 {
		return false, ""
	} else {
		return true, hash
	}

}
