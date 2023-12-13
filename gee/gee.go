package gee

import (
	"log"
	"net/http"
)

/*
定义gee使用的处理器函数
定义HandlerFunc,提供给框架用户的，用来定义路由映射的处理方法。
在Engine中，添加了一张路由映射表router，key由请求方法和静态路由地址构成。例如 GET-/、GET-/hello
这样针对相同的路由，如果请求方法不同，可以映射不同处理方法（Handler）,value是用户映射的处理方法。
*/
type HandlerFunc func(c *Context)

// Engine 创建引擎 Engine 实现 ServeHTTP接口
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有路由组
}

// RouterGroup 路由组
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持嵌套
	engine      *Engine       // 所有组共享这个引擎
}

// New 构建 gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	return newGroup
}

/*
将请求添加到路由映射表中router
参数：

	method:请求的方法
	comp: URL
	handler:处理器函数
*/
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

/*
GET
GET请求的方法
参数:

	pattern: URL
	handler: 处理器函数
*/
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

/*
POST
POST请求的方法
参数:

	pattern: URL
	handler: 处理器函数
*/
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

// Run 定义run方法开启http服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

/*
定义Engine结构体的ServeHTTP方法
ServeHTTP 方法的作用
解析请求的路径，查找路由映射表。
如果查找到，就执行注册的处理方法，
如果没有查找到，就返回 404 NOT FOUND。
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.Handle(c)
}
