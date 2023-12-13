package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {

		// 开始时间
		t := time.Now()
		// 进程请求
		c.Next()

		//计算决定时间
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
