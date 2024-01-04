package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (j *JSON) Render(w http.ResponseWriter) error {
	return WriteJSON(w, j.Data)
}

func WriteJSON(w http.ResponseWriter, data any) error {
	writeContentType(w, jsonContentType)
	err := json.NewEncoder(w).Encode(data)
	return err
}
func (j *JSON) WriteHeader(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
