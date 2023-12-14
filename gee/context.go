package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context 上下文结构体
type Context struct {
	// 起始对象
	Writer http.ResponseWriter
	Req    *http.Request

	// 请求信息
	Path   string
	Method string
	Params map[string]string

	// 响应信息
	StatusCode int

	// 中间件
	handlers []HandlerFunc
	index    int

	// 引擎指针
	engine *Engine
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

/*
newContext
创建上下文
参数:

	w: http.ResponseWriter
	req: *http.Request

返回 *Context
*/
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

/*
PostForm
获取post请求中的参数
*/
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置响应的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code) // 设置响应状态码
}

/*
SetHeader
设置响应头信息
*/
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

/*
String
以字符串的形式，返回数据
*/
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

/*
JSON
以JSON形式返回数据
*/
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

/*
Data
返回数据
*/
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

/*
HTML
以HTML格式将数据返回
*/
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, msg string) {
	c.String(code, msg)
}
