package ddgo

import (
	"github.com/a754962942/ddgo/render"
	"github.com/a754962942/ddgo/utils"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	engine     *Engine
	queryCache url.Values
}

func (c *Context) Param(s string, ss any) string {
	s2 := ss.(string)
	return "[" + s + " " + s2 + "]"
}

func (c *Context) HTML(statCode int, body string) {
	err := c.Render(statCode, &render.HTML{Data: body, IsTemplate: false})
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

func (c *Context) Template(statCode int, name string, data interface{}) {
	err := c.Render(statCode, &render.HTML{
		IsTemplate: true,
		Name:       name,
		Data:       data,
		Template:   c.engine.HTMLRender.Template,
	})
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) JSON(statCode int, data any) {
	err := c.Render(statCode, &render.JSON{Data: data})
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) String(statCode int, format string, value ...any) {
	err := c.Render(statCode, &render.String{Format: format, Data: value})
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) Render(code int, r render.Render) error {
	if _, ok := r.(*render.Redirect); !ok {
		c.W.WriteHeader(code)
	}
	err := r.Render(c.W)
	return err
}

func (c *Context) XML(statCode int, data any) {
	err := c.Render(statCode, &render.XML{Data: data})
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) File(filename string) {
	http.ServeFile(c.W, c.R, filename)
}

func (c *Context) FileAttachment(filepath, filename string) {
	if utils.IsASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

// FileFromFS
//
//	@Description: filepath为fs下对应路径
//	@receiver c
//	@param filepath
//	@param fs
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) Redirect(statCode int, url string) {
	redirect := render.Redirect{
		Code:    statCode,
		Url:     url,
		Request: c.R,
	}
	err := c.Render(statCode, &redirect)
	if err != nil {
		log.Println(err)
	}
}
func (c *Context) DefaultQuery(key, defaultValue string) string {
	array, ok := c.GetQueryArray(key)
	if !ok {
		return defaultValue
	}
	return array[0]
}
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}
func (c *Context) QueryArray(key string) (values []string) {
	c.initQueryCache()
	values, _ = c.queryCache[key]
	return
}
func (c *Context) GetQuery(key string) string {
	c.initQueryCache()
	return c.queryCache.Get(key)
}
func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.R != nil {
			c.queryCache = c.R.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}
