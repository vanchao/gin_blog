package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"
	sq "github.com/Masterminds/squirrel"
)

type Post struct {
	Id        int    `json:"id" form:"id"`
	TagId     int    `json:"tag_id" form:"tag_id"`
	AuthorId  int    `json:"author_id" form:"author_id"`
	UpdatedTime sql.NullString      `json:"updated_time" form:"updated_time"`
	CreatedTime string      `json:"created_time" form:"created_time"`
	AuthorName string      `json:"author_name" form:"author_name"`
	Title string `json:"title" form:"title"`
	TagName string `json:"tag_name" form:"tag_name"`
	Description string `json:"description" form:"description"`
	Content string `json:"content" form:"content"`
}


func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT"}
	config.AllowOrigins = []string{"http://127.0.0.1:5550/"}
	r.Use(cors.New(cors.Config{
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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/post_list", func(c *gin.Context) {
		var offset, limit uint64

		offset	 = 0

		limit 	 = 5

		sql_string, args, err := sq.Select("*").From("posts").Offset(offset).Limit(limit).ToSql()

		if(nil != err) {
			e := fmt.Sprintf("Error found40: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": e,
			})
		}

		sqlConnString := getConnString()
		db, err := sql.Open("mysql", sqlConnString)
		if err != nil {
			e := fmt.Sprintf("Error found50: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": e,
			})
		}

		defer db.Close()
		rows, err := db.Query(sql_string)
		if err != nil {
			e := fmt.Sprintf("Error found60: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": e,
			})
		}

		posts := make([]Post, 0)

		for rows.Next() {
			var post Post
			rows.Scan(&post.Id, &post.TagId, &post.AuthorId, &post.UpdatedTime, &post.CreatedTime, &post.AuthorName, &post.Title, &post.TagName, &post.Description, &post.Content)
			posts = append(posts, post)
		}

		_ = args

		if(nil != err) {
			e := fmt.Sprintf("Error found70: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": e,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"result": posts,
		})

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getConnString() string {
	viper.SetConfigName("database")
	viper.SetConfigType("json")
	viper.AddConfigPath("./main/config/")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	uname := viper.Get("user_name")
	pwd := viper.Get("password")
	protocol := viper.Get("protocol")
	db := viper.Get("database")
	str := fmt.Sprintf("%s:%s@%s/%s", uname, pwd, protocol, db)
	return str
}
