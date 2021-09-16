package main

import (
	"go_toy"
	"log"
	"net/http"
	"time"
)

func onlyForV2() go_toy.HandlerFunc {
	return func(c *go_toy.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.Fail(500, "Internal Server Error")
		c.StatusCode = 500
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := go_toy.New()
	r.Use(go_toy.Logger())
	r.Use(go_toy.Recovery())
	r.GET("/", func(c *go_toy.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Go Toy</h1>")
	})

	r.GET("/panic", func(c *go_toy.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/index", func(c *go_toy.Context){
			c.HTML(http.StatusOK, "hello index")
		})

		v1.GET("/hello/:name", func(c *go_toy.Context){
			c.JSON(http.StatusOK, go_toy.H{"hello": c.Param("name")})
		})
	}

	r.GET("/hello", func(c *go_toy.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *go_toy.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *go_toy.Context) {
		c.JSON(http.StatusOK, go_toy.H{"filepath": c.Param("filepath")})
	})
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	v2.GET("/test", func(c *go_toy.Context) {
		c.JSON(http.StatusOK, go_toy.H{"test": "ok"})
	})


	r.Run(":9999")
}