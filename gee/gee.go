/*
 * @Autor: qing fu liu
 * @Email: 1805003304@qq.com
 * @Github: https://github.com/qingfuliu
 * @Date: 2022-03-21 09:29:46
 * @LastEditors: qingfu liu
 * @LastEditTime: 2022-03-21 21:18:19
 * @FilePath: \golang\Gee\gee\gee.go
 * @Description:
 */
package gee

import (
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

//路由分组
type routeGroup struct {
	engine     *Engine
	parent     *routeGroup
	children   map[string]*routeGroup
	prefix     string
	middleware []HandleFunc
}

func (r *routeGroup) Group(path string) *routeGroup {
	parts := parserPath(path)
	temp := r
	for len(parts) > 0 {
		value, ok := temp.children[parts[0]]
		if !ok {
			temp.addGroup(parts[0])
			value = temp.children[parts[0]]
		}
		temp = value
		parts = parts[1:]
	}
	return temp
}

func (r *routeGroup) addGroup(path string) {
	parts := parserPath(path)
	temp := r
	for len(parts) > 0 {
		value, ok := temp.children[parts[0]]
		if !ok {
			value = &routeGroup{
				engine:   temp.engine,
				parent:   temp,
				children: make(map[string]*routeGroup),
				prefix:   temp.prefix + "/" + parts[0]}

			temp.children[parts[0]] = value
		}
		temp = value
		parts = parts[1:]
	}
}

type Engine struct {
	*routeGroup
	router_ *router
}

func New() *Engine {
	engine_ := &Engine{router_: newRouter()}
	engine_.routeGroup = &routeGroup{engine: engine_, parent: nil, children: make(map[string]*routeGroup), prefix: "/", middleware: make([]HandleFunc, 0)}
	return engine_
}

func (engine *Engine) Use(group string, middlewares ...HandleFunc) {
	group_ := engine.Group(group)
	group_.Use(middlewares...)
}

func (g *routeGroup) Use(middlewares ...HandleFunc) {
	g.middleware = append(g.middleware, middlewares...)
}

func (engine *Engine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	c := newContext(response, request)
	for key, value := range engine.children {
		if strings.HasPrefix(request.URL.Path, key) {
			c.handles = append(c.handles, value.middleware...)
		}
	}
	engine.router_.handle(c)
}

func (r *routeGroup) addRoute(method, path string, f HandleFunc) {
	r.engine.router_.addRouter(method, r.prefix+"/"+path, f)
}

func (r *routeGroup) deleteRoute(method, path string) {
	r.engine.router_.deleteRoute(method, r.prefix+"/"+path)
}

func (r *routeGroup) Get(pattern string, handle HandleFunc) {
	r.engine.router_.addRouter("GET", r.prefix+"/"+pattern, handle)
}

func (r *routeGroup) Post(pattern string, handle HandleFunc) {
	r.engine.router_.addRouter("POST", r.prefix+"/"+pattern, handle)
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}
