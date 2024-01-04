package ddgo

import (
	"fmt"
	"github.com/a754962942/ddgo/render"
	"github.com/a754962942/ddgo/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const ANY = "ANY"

type HandlerFunc func(context *Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type routerGroup struct {
	name               string
	handlerFuncMap     map[string]map[string]HandlerFunc
	middlewaresFuncMap map[string]map[string][]MiddlewareFunc
	//冗余
	handlerMethodMap map[string][]string
	treeNode         *treeNode
	middlewares      []MiddlewareFunc
}

func (r *routerGroup) PreHandle(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(routerName string, method string, h HandlerFunc, ctx *Context) {
	//middlewares
	if r.middlewares != nil {
		for _, middlewareFunc := range r.middlewares {
			h = middlewareFunc(h)
		}
	}
	middlewareFuncs := r.middlewaresFuncMap[routerName][method]
	if middlewareFuncs != nil {
		for _, middlewareFunc := range middlewareFuncs {
			h = middlewareFunc(h)
		}
	}
	h(ctx)
}

type router struct {
	routerGroups []*routerGroup
}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	left := strings.TrimLeft(name, "/")
	_, ok := r.handlerFuncMap["/"+left]
	if !ok {
		r.handlerFuncMap["/"+left] = make(map[string]HandlerFunc)
		r.middlewaresFuncMap["/"+left] = make(map[string][]MiddlewareFunc)
	}
	_, ok = r.handlerFuncMap["/"+left][method]
	if ok {
		panic(method + " method already exist")
	}
	r.handlerFuncMap["/"+left][method] = handlerFunc
	r.middlewaresFuncMap["/"+left][method] = middlewareFunc
	r.treeNode.Put("/" + left)
}

func (r *routerGroup) GET(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) POST(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handlerFunc, middlewareFunc...)
}

func (r *routerGroup) ANY(name string, handlerFunc HandlerFunc, middlewareFunc ...MiddlewareFunc) {
	r.handle(name, ANY, handlerFunc, middlewareFunc...)
}

func (r *router) Group(name string) *routerGroup {
	left := strings.TrimLeft(name, "/")
	rg := &routerGroup{
		name:             left,
		handlerFuncMap:   make(map[string]map[string]HandlerFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode: &treeNode{
			name:     "/" + left,
			children: make([]*treeNode, 0),
		},
		middlewaresFuncMap: make(map[string]map[string][]MiddlewareFunc),
	}
	r.routerGroups = append(r.routerGroups, rg)
	return rg
}

type Engine struct {
	funcMap    template.FuncMap
	HTMLRender render.HTMLRender
	router
	pool sync.Pool
}

func New() *Engine {
	engine := &Engine{
		router: router{routerGroups: make([]*routerGroup, 0)},
	}
	engine.pool.New = func() any {
		return engine.allocateContext()
	}
	return engine
}

func (e *Engine) allocateContext() any {
	return &Context{engine: e}
}
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
	e.HTMLRender = render.HTMLRender{Template: t}
}

func (e *Engine) SetHtmlTemplate(t *template.Template) {
	e.HTMLRender = render.HTMLRender{Template: t}
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.W = w
	ctx.R = r
	e.httpRequestHandler(ctx)
	e.pool.Put(ctx)
}

func (e *Engine) Run(port string) {
	server := http.Server{
		Addr:         port,
		Handler:      e,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (e *Engine) httpRequestHandler(ctx *Context) {
	method := ctx.R.Method
	//fmt.Println("method:", method)
	for _, group := range e.routerGroups {
		groupName := utils.SubStrAndTrim(ctx.R.URL.Path, "/"+group.name)
		node := group.treeNode.Get(groupName)
		if node != nil && node.isEnd {

			if handle, ok := group.handlerFuncMap[node.routerName][ANY]; ok {
				group.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			if handle, ok := group.handlerFuncMap[node.routerName][method]; ok {
				group.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			ctx.W.WriteHeader(http.StatusMethodNotAllowed)
			_, err := fmt.Fprintln(ctx.W, ctx.R.RequestURI+" "+method+" Not ALLOW")
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	ctx.W.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintln(ctx.W, ctx.R.RequestURI+" Not Found")
	if err != nil {
		log.Println(err)
	}
}
