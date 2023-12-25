package ddgo

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HandlerFunc func(context Context)

const ANY = "ANY"

type router struct {
	routerGroups []*routerGroup
}

type routerGroup struct {
	name           string
	handlerFuncMap map[string]map[string]HandlerFunc
	//冗余
	handlerMethodMap map[string][]string
	treeNode         *treeNode
}

func (r *routerGroup) handle(name string, method string, handlerFunc HandlerFunc) {
	left := strings.TrimLeft(name, "/")
	_, ok := r.handlerFuncMap[left]
	if !ok {
		r.handlerFuncMap[left] = make(map[string]HandlerFunc)
	}
	_, ok = r.handlerFuncMap[left][method]
	if ok {
		panic(method + " method already exist")
	}
	r.handlerFuncMap[left][method] = handlerFunc
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
	method := r.Method

	//fmt.Println("method:", method)
	for _, group := range e.routerGroups {
		groupName := SubStrAndTrim(r.RequestURI, "/"+group.name)
		node := group.treeNode.Get(groupName)
		if node != nil {
			context := Context{
				W: w,
				R: r,
			}
			if handle, ok := group.handlerFuncMap[groupName][ANY]; ok {
				handle(context)
				return
			}
			if handle, ok := group.handlerFuncMap[groupName][method]; ok {
				handle(context)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, r.RequestURI+" "+method+" Not ALLOW")
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, r.RequestURI+" Not Found")
	return
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
