package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

type ListChecks struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func InitMysql() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/GORM01?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error creating database")
		return
	}
	return DB.DB().Ping()
}
func main() {
	//创建数据库
	//连接数据库
	err := InitMysql()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	DB.AutoMigrate(&ListChecks{})
	defer DB.Close()

	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLFiles("./templates/index.html")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1Group := r.Group("v1")
	{
		//添加一条待办事项
		v1Group.POST("/todo", func(c *gin.Context) {
			var Todo ListChecks
			c.BindJSON(&Todo)
			err = DB.Create(&Todo).Error
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, Todo)
				fmt.Println(&Todo)
			}
		})

		//查询所有代办
		v1Group.GET("/todo", func(c *gin.Context) {
			var AllList []ListChecks
			err = DB.Find(&AllList).Error
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, AllList)
			}
		})
		//更新一条
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			//bool默认为false
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "不存在该ID"})
			}
			var findOne ListChecks
			err = DB.Where("id=?", id).First(&findOne).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error()})
			}

			if err = DB.Save(&findOne).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, findOne)
			}
		})
		//删除一条
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			// 前端传入删除id
			// 后端处理id

			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				fmt.Println(err.Error())
			}
			err = DB.Where("id=?", id).Delete(ListChecks{}).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "delete"})
			}
		})
	}

	r.Run(":9000")

}
