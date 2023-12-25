package ddgo

import (
	"fmt"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) String(msg string) {
	fmt.Fprintln(c.W, msg)
}
