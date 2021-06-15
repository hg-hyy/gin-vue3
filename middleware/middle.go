package middleware

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")
		// Host := c.Request.Header.Get("Host")
		// fmt.Println(Host, origin)
		if len(origin) == 0 {
			origin = c.Request.Header.Get("Origin")
			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		} else {
			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,authtoken")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST,DELETE")
		// c.Writer.Header().Set("Content-Type", "application/json;text/html; charset=utf-8")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		tm := time.Now()
		path := c.Request.URL.Path
		rw := c.Request.URL.RawPath
		filer, err := os.OpenFile("logs/pt.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Println(err)
		}

		gin.DefaultWriter = io.MultiWriter(filer)
		log.New(filer, "", log.Lshortfile)
		log.SetOutput(filer)
		c.Next()

		if path == "" {
			path = path + rw
		}
		latency := time.Since(tm)
		method := c.Request.Method
		status := c.Writer.Status()
		log.Printf("[%s] [%s] [%d] [%s]", method, path, status, latency)

	}
}
