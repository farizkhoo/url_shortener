package main

import (
	"log"

	"github.com/farizkhoo/url_shortener/cmd/api/app"
)

func main() {
	pg := app.NewPostgres()
	if err := pg.Connect(); err != nil {
		log.Fatalln(err)
	}

	uh := app.NewURLHandler(pg)
	r := app.SetupRouter(uh)

	r.Run()
}
