package main

import (
	"github.com/gin-gonic/gin"
)

var pg *Postgres

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/shorten_url", createURLHandler)
	r.POST("/shorten_url/", createURLHandler)
	r.GET("/shorten_url/:id", getURLHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return r
}

func main() {
	pg = NewPostgres()

	r := setupRouter()

	r.Run()
}
