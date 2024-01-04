package render

import (
	"errors"
	"fmt"
	"net/http"
)

type Redirect struct {
	Code    int
	Url     string
	Request *http.Request
}

func (r *Redirect) Render(w http.ResponseWriter) error {
	//r.WriteHeader(w)
	if (r.Code < http.StatusMultipleChoices ||
		r.Code > http.StatusPermanentRedirect) && r.Code != http.StatusCreated {
		return errors.New(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Url, r.Code)
	return nil
}
func (r *Redirect) WriteHeader(w http.ResponseWriter) {

}
