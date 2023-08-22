package gee

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 定义个函数体的别名为HandlerFunc
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middleWares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	//继承routerGroup的方法和变量
	*RouterGroup
	router       *Router
	groups       []*RouterGroup   //store all groups
	htmlTemplate *template.HTML   //for html render
	funMap       template.FuncMap // for html render
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger())
	engine.Use(Recovery())
	return engine
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {

	newGroup := &RouterGroup{
		engine: group.engine,
		prefix: group.prefix + prefix,
		parent: group,
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {

	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleWares []HandlerFunc
	for _, group := range engine.groups {
		//得到属于你自己组的中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}
	c := NewContext(w, req)
	//把得到的中间件方法放到context里面
	c.handlers = middleWares
	engine.router.handle(c)
}
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use 定义Use函数，将中间件应用到某个 Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middleWares = append(group.middleWares, middlewares...)
}
