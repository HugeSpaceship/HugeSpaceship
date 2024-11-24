package utils

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type TmplMap map[string]any

const xmlMarshalError = "Failed to marshal xml"

func HttpLog(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	slog.Debug("HTTP Log", "status", status, "message", message)
	if err != nil {
		slog.Error("failed to write to ResponseWriter")
	}
}

func HttpLogf(w http.ResponseWriter, status int, format string, a ...any) {
	HttpLog(w, status, fmt.Sprintf(format, a...))
}

func GetContextValue[T any](ctx context.Context, key string) T {
	return ctx.Value(key).(T)
}

func XMLMarshal(w http.ResponseWriter, o any) error {
	w.Header().Set("Content-Type", "text/xml")
	slotBytes, err := xml.Marshal(o)
	if err != nil {
		return err
	}
	_, err = w.Write(slotBytes)
	if err != nil {
		return err
	}
	return nil
}

func XMLUnmarshal[T any](r *http.Request) (*T, error) {
	out := new(T)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(data, out)
	return out, err
}
