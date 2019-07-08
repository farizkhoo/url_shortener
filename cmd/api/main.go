package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/shorten_url", createURLHandler)
	r.POST("/shorten_url/", createURLHandler)
	r.GET("/shorten_url/:id", getURLHandler)

	return r
}

func main() {
	var err error
	db, err = NewPostgres().Connect()
	if err != nil {
		log.Fatalln(err)
	}

	r := setupRouter()

	r.Run()
}
