package audit

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteJSON serialises the audit Record as indented JSON to w.
func WriteJSON(w io.Writer, r Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("audit: failed to encode JSON: %w", err)
	}
	return nil
}

// ReadJSON deserialises a Record from JSON in r.
func ReadJSON(r io.Reader) (Record, error) {
	var rec Record
	if err := json.NewDecoder(r).Decode(&rec); err != nil {
		return Record{}, fmt.Errorf("audit: failed to decode JSON: %w", err)
	}
	return rec, nil
}
