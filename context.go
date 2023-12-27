package ddgo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	engine *Engine
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

func (c *Context) HTML(statCode int, body string) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(statCode)
	_, err := c.W.Write([]byte(body))
	if err != nil {
		log.Println(err)
	}
}
func (c *Context) HTMLTemplate(statCode int, name string, data any, filenames ...string) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(statCode)
	t := template.New(name)
	files, err := t.ParseFiles(filenames...)
	if err != nil {
		log.Println(err)
	}
	err = files.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}
func (c *Context) HTMLTemplateGlob(statCode int, name string, data any, parttern string) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(statCode)
	t := template.New(name)
	files, err := t.ParseGlob(parttern)
	if err != nil {
		log.Println(err)
	}
	err = files.Execute(c.W, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) Template(name string, data interface{}) {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(http.StatusOK)
	err := c.engine.HTMLRender.Template.ExecuteTemplate(c.W, name, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) JSON(statCode int, data any) {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(statCode)
	marshal, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	_, err = c.W.Write(marshal)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) XML(statCode int, data any) {
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(statCode)
	err := xml.NewEncoder(c.W).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
