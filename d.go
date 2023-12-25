package ddgo

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const ANY = "ANY"

type HandlerFunc func(context Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type routerGroup struct {
	name           string
	handlerFuncMap map[string]map[string]HandlerFunc
	//冗余
	handlerMethodMap map[string][]string
	treeNode         *treeNode
	middlewares      []MiddlewareFunc
}

func (r *routerGroup) PreHandle(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(h HandlerFunc, ctx Context) {
	if r.middlewares != nil {
		for _, middlewareFunc := range r.middlewares {
			h = middlewareFunc(h)
		}
	}
	h(ctx)
}

type router struct {
	routerGroups []*routerGroup
}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc) {
	left := strings.TrimLeft(name, "/")
	_, ok := r.handlerFuncMap["/"+left]
	if !ok {
		r.handlerFuncMap["/"+left] = make(map[string]HandlerFunc)
	}
	_, ok = r.handlerFuncMap["/"+left][method]
	if ok {
		panic(method + " method already exist")
	}
	r.handlerFuncMap["/"+left][method] = handlerFunc
	//冗余
	r.handlerMethodMap[method] = append(r.handlerMethodMap[method], left)
	r.treeNode.Put("/" + left)
}

func (r *routerGroup) GET(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodGet, handlerFunc)
}

func (r *routerGroup) POST(name string, handlerFunc HandlerFunc) {
	r.handle(name, http.MethodPost, handlerFunc)
}

func (r *routerGroup) ANY(name string, handlerFunc HandlerFunc) {
	r.handle(name, ANY, handlerFunc)
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
	}
	r.routerGroups = append(r.routerGroups, rg)
	return rg
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{routerGroups: make([]*routerGroup, 0)},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.httpRequestHandler(w, r)
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

func (e *Engine) httpRequestHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	//fmt.Println("method:", method)
	for _, group := range e.routerGroups {
		groupName := SubStrAndTrim(r.RequestURI, "/"+group.name)
		node := group.treeNode.Get(groupName)
		if node != nil && node.isEnd {
			context := Context{
				W: w,
				R: r,
			}
			if handle, ok := group.handlerFuncMap[node.routerName][ANY]; ok {
				group.methodHandle(handle, context)
				return
			}
			if handle, ok := group.handlerFuncMap[node.routerName][method]; ok {
				group.methodHandle(handle, context)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err := fmt.Fprintln(w, r.RequestURI+" "+method+" Not ALLOW")
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintln(w, r.RequestURI+" Not Found")
	if err != nil {
		log.Println(err)
	}
}
