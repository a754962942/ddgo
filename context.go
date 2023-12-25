package ddgo

import (
	"fmt"
	"log"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) String(msg ...string) {
	_, err := fmt.Fprintln(c.W, msg)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) Param(s string, ss any) string {
	s2 := ss.(string)
	return "[" + s + " " + s2 + "]"
}
