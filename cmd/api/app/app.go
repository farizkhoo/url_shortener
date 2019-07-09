package app

import "github.com/gin-gonic/gin"

// URL holds the uuid and the name of the Long
type URL struct {
	ID      string
	LongURL string
}

// SetupRouter sets up routes and returns gin engine
func SetupRouter(u *URLHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/shorten_url/:id", u.GetURL())
	r.POST("/shorten_url/", u.CreateURL())

	return r
}
