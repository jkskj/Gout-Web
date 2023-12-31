package main

import (
	"Gout"
	"log"
	"net/http"
	"time"
)

func onlyForV2() Gout.HandlerFunc {
	return func(c *Gout.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := Gout.Default()
	r.GET("/", func(c *Gout.Context) {
		c.JSON(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello/:name", func(c *Gout.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	r.GET("/hello/*file", func(c *Gout.Context) {
		c.JSON(http.StatusOK, Gout.H{"file": c.Param("file")})
	})
	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *Gout.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run()
}
