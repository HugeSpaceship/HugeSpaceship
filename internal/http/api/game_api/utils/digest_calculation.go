package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log/slog"
)

func CalculateDigest(path, authCookie, digestKey string, body []byte, excludeBody bool) string {
	sum := sha1.New()

	if !excludeBody {
		sum.Write(body)
	}

	if _, err := io.WriteString(sum, authCookie); err != nil {
		slog.Debug("Failed to write auth cookie", "error", err)
	}

	if _, err := io.WriteString(sum, path); err != nil {
		slog.Debug("Failed to write path", "error", err)
	}

	if _, err := io.WriteString(sum, digestKey); err != nil {
		slog.Debug("Failed to write digest key", "error", err)
	}

	return hex.EncodeToString(sum.Sum(nil))
}
