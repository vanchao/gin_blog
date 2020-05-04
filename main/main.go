package main

import (
	_ "encoding/json"
	"gin_blog/main/frontend"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	_ "github.com/go-sql-driver/mysql"
)


func main() {
	r := gin.Default()
	v1 := r.Group("/frontend")
	{
		config := cors.DefaultConfig()
		config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT"}
		config.AllowOrigins = []string{"http://127.0.0.1:5550/"}
		v1.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://127.0.0.1:5500"},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				//return origin == "https://github.com"
				return true
			},
			MaxAge: 12 * time.Hour,
		}))
		v1.GET("/post_list", frontend.Post_list)
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

