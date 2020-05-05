package frontend

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"math"
	"net/http"
	"strconv"
)

const (
	IS_DISPLAY  = 10
	NOT_DISPLAY = 20
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

var (
	total_count, offset, limit,  first_page, next_page, prev_page  uint64
)


func Post_list(c *gin.Context) {

	c_page := c.Query("current_page")

	current_page, err := strconv.ParseUint(c_page, 10, 64)
	if err != nil {
		e := fmt.Sprintf("Error found30: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}

	first_page = 1

	if 0 == current_page {
		current_page = first_page
	}

	limit 	 = 5

	offset	 = (current_page - 1) * limit


	select_string, args, err := sq.Select("*").From("posts").OrderBy("created_time DESC").Offset(offset).Limit(limit).ToSql()

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

	count_string, args, err := sq.Select("COUNT(id) as COUNT").From("posts").ToSql()
	if(nil != err) {
		e := fmt.Sprintf("Error found41: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}


	cnt_rows, err := db.Query(count_string)

	var count int

	for cnt_rows.Next() {
		if err := cnt_rows.Scan(&count); err != nil {
			e := fmt.Sprintf("Error found42: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": e,
			})
		}
	}

	if(nil != err) {
		e := fmt.Sprintf("Error found43: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}

	rows, err := db.Query(select_string)
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

	total_count = uint64(count)

	last_page := math.Ceil( (float64(total_count * 100) / float64(limit)) / 100  )

	next_page = current_page + 1

	prev_page = current_page - 1

	is_prev_display := IS_DISPLAY

	is_next_display := IS_DISPLAY

	if current_page >= uint64(last_page)  {
		current_page = uint64(last_page)
		is_next_display = NOT_DISPLAY
	}

	if next_page > uint64(last_page) {
		next_page = uint64(last_page)
	}

	if prev_page < first_page {
		prev_page = first_page
	}

	if current_page <= prev_page {
		current_page = prev_page
		is_prev_display = NOT_DISPLAY
	}

	c.JSON(http.StatusOK, gin.H{
		"result": posts,
		"current_page": current_page,
		"next_page": next_page,
		"prev_page": prev_page,
		"is_next_display": is_next_display,
		"is_prev_display": is_prev_display,
	})

}

func Post_info(c *gin.Context) {

	//Id        int    `json:"id" form:"id"`
	//TagId     int    `json:"tag_id" form:"tag_id"`
	//AuthorId  int    `json:"author_id" form:"author_id"`
	//UpdatedTime sql.NullString      `json:"updated_time" form:"updated_time"`
	//CreatedTime string      `json:"created_time" form:"created_time"`
	//AuthorName string      `json:"author_name" form:"author_name"`
	//Title string `json:"title" form:"title"`
	//TagName string `json:"tag_name" form:"tag_name"`
	//Description string `json:"description" form:"description"`
	//Content string `json:"content" form:"content"`


	id := c.Query("post_id")

	var tag_id, author_id int

	var created_time, author_name, title, tag_name,  description, content string

	post_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		e := fmt.Sprintf("Error found30: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}

	select_active := sq.Select("tag_id, author_id, created_time, author_name, title, tag_name, description, content").From("posts").Where(sq.Eq{"id": post_id})

	sqlConnString := getConnString()
	db, err := sql.Open("mysql", sqlConnString)
	if err != nil {
		e := fmt.Sprintf("Error found50: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}

	defer db.Close()

	rows, err := select_active.RunWith(db).Query()
	if err != nil {
		e := fmt.Sprintf("Error found60: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_msg": e,
		})
	}

	for rows.Next() {
		if err := rows.Scan(&tag_id, &author_id, &created_time, &author_name, &title, &tag_name, &description, &content); err != nil {
			e := fmt.Sprintf("Error found42: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_msg": e,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id": post_id,
		"created_time": created_time,
		"author_name": author_name,
		"content": content,
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

