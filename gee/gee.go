package gee

import (
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
	group  []*RouterGroup // 存储所有路由组
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持嵌套
	engine      *Engine       // 所有组共享这个引擎
}

// New 构建 gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

/*
将请求添加到路由映射表中router
参数：

	method:请求的方法
	pattern: URL
	handler:处理器函数
*/
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRouter(method, pattern, handler)
}

/*
GET
GET请求的方法
参数:

	pattern: URL
	handler: 处理器函数
*/
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

/*
POST
POST请求的方法
参数:

	pattern: URL
	handler: 处理器函数
*/
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
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
