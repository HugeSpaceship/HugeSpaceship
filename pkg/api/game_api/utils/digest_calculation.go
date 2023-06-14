package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"io"
)

func CalculateDigest(path, authCookie, digestKey string, body []byte, excludeBody bool) string {
	sum := sha1.New()

	if !excludeBody {
		sum.Write(body)
	}

	if _, err := io.WriteString(sum, authCookie); err != nil {
		log.Debug().Err(err).Msg("Failed to write auth cookie")
	}

	if _, err := io.WriteString(sum, path); err != nil {
		log.Debug().Err(err).Msg("Failed to write path")
	}

	if _, err := io.WriteString(sum, digestKey); err != nil {
		log.Debug().Err(err).Msg("Failed to write digest key")
	}

	return hex.EncodeToString(sum.Sum(nil))
}
