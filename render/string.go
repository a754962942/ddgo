package render

import (
	"fmt"
	"github.com/a754962942/ddgo/utils"
	"net/http"
)

type String struct {
	Format string
	Data   []any
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func (s *String) Render(w http.ResponseWriter) error {
	return WriteString(w, s.Format, s.Data)
}

func WriteString(w http.ResponseWriter, format string, data []any) error {
	writeContentType(w, plainContentType)
	if len(data) > 0 {
		_, err := fmt.Fprintf(w, format, data...)
		return err
	}
	_, err := w.Write(utils.StringToBytes(format))
	return err
}

func (s *String) WriteHeader(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}
