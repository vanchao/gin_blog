package frontend

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	sq "github.com/Masterminds/squirrel"
	"strconv"
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


func Post_list(c *gin.Context) {

	c_page := c.Query("current_page")

	current_page, err := strconv.ParseUint(c_page, 10, 64)
	if err != nil {
		e := fmt.Sprintf("Error found30: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}

	if 0 == current_page {
		current_page = 1
	}

	var offset, limit uint64

	limit 	 = 5

	offset	 = (current_page - 1) * limit

	sql_string, args, err := sq.Select("*").From("posts").OrderBy("created_time DESC").Offset(offset).Limit(limit).ToSql()

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
		"current_page": current_page,
	})

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

