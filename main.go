package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// 开始时间
		t := time.Now()
		// 如果错误，就返回错误
		c.Fail(500, "Internal Server Error")
		// 计算最终时间
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()

	r.Use(gee.Logger())
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<H1>Hello Gee</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Helo Liu Zhenwei</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "Hello %s,you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s,you're at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s,you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"filepath": c.Param("filepath"),
		})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8000")
}
