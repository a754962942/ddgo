package render

import "net/http"

type Render interface {
	Render(w http.ResponseWriter) error
	WriteHeader(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, conentType []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = conentType
	}
}
