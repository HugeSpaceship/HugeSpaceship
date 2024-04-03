package utils

import (
	"encoding/xml"
	"net/http"
)

// LBPXML contains the given interface object.
type LBPXML struct {
	Data any
}

var xmlContentType = []string{"text/xml"}

// Render (LBPXML) encodes the given interface object and writes data with custom ContentType.
func (r LBPXML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

// WriteContentType (LBPXML) writes LBPXML ContentType for response.
func (r LBPXML) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = xmlContentType
	}
}
