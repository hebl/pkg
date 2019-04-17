package httputils

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// WriteJSON writes the value v to the http response stream as json with standard json encoding.
func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	var err error
	if b, ok := v.(string); ok {
		_, err = w.Write([]byte(b))
	} else {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		return enc.Encode(v)
	}
	return err
}

// WriteXML writes the value v to the http response stream as json with standard json encoding.
func WriteXML(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	var err error
	if b, ok := v.(string); ok {
		_, err = w.Write([]byte(b))
	} else {
		enc := xml.NewEncoder(w)
		return enc.Encode(v)
	}
	return err
}
