package render

import (
	"github.com/a754962942/ddgo/utils"
	"html/template"
	"net/http"
)

type HTMLRender struct {
	Template *template.Template
}
type HTML struct {
	Data       any
	Name       string
	Template   *template.Template
	IsTemplate bool
}

func (h *HTML) Render(w http.ResponseWriter) error {
	h.WriteHeader(w)
	if h.IsTemplate {
		return h.Template.ExecuteTemplate(w, h.Name, h.Data)
	}
	_, err := w.Write(utils.StringToBytes(h.Data.(string)))
	return err
}
func (h *HTML) WriteHeader(w http.ResponseWriter) {
	writeContentType(w, []string{"text/html;charset=utf-8"})
}
