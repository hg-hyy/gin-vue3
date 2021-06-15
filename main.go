package main

import (
	"flag"
	"log"
	"net/http"
	"pt/api"
	"pt/middleware"

	"github.com/gin-gonic/gin"
)

var (
	Addr = flag.String("addr", "127.0.0.1:8002", "http server host")
)

func main() {
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	var api api.Api
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	r.Static("/assets", "./assets")
	r.StaticFS("/static", http.Dir("./static"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.LoadHTMLGlob("templates/**/*") //模板查找一定要放在下面

	api.Routers(r)
	s := http.Server{
		Addr:    *Addr,
		Handler: r,
	}
	log.Fatal(s.ListenAndServe())
}
