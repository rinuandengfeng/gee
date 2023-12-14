package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
	router        *router
	groups        []*RouterGroup     // 存储所有路由组
	htmlTemplates *template.Template // html渲染  将所有模版加载进内存
	funcMap       template.FuncMap   // html渲染  自定义模版函数
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

// Use 是定义将中间件添加到路由组中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// 如果收到的请求路由前缀在组前缀中，就将其中间件添加到中间件中
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.Handle(c)

}

// createStaticHandler 创建静态处理器
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否存在，或者 是否有权限
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 静态文件服务
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册GET处理器
	group.GET(urlPattern, handler)
}

// SetFuncMap 设置模版函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载模版方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
