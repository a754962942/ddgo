package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data any
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

func (x *XML) Render(w http.ResponseWriter) error {
	return WriteXML(w, x.Data)
}

func WriteXML(w http.ResponseWriter, data any) error {
	writeContentType(w, xmlContentType)
	return xml.NewEncoder(w).Encode(data)
}
func (x *XML) WriteHeader(w http.ResponseWriter) {
	writeContentType(w, xmlContentType)
}
