package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Blog struct {
	*gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Table    string
}

func getDBConfig() DBConfig {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Table:    os.Getenv("DB_NAME"),
	}
}

func connectDB() (*gorm.DB, error) {
	config := getDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", config.User, config.Password, config.Host, config.Port, config.Table)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db, err
}

func errorDB(db *gorm.DB, c *gin.Context) bool {
	if db.Error != nil {
		log.Printf("Error todos: %v", db.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return true
	}
	return false
}

func listener(r *gin.Engine, db *gorm.DB) {
	r.GET("/blog/delete", func(c *gin.Context) {
		id, _ := c.GetQuery("id")
		result := db.Delete(&Blog{}, id)
		if errorDB(result, c) {
			return
		}
	})
	r.POST("/blog/update", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		title := c.PostForm("title")
		content := c.PostForm("content")
		var blog Blog
		result := db.Where("id = ?", id).Take(&blog)
		if errorDB(result, c) {
			return
		}
		blog.Title = title
		blog.Content = content
		result = db.Save(&blog)
		if errorDB(result, c) {
			return
		}
	})
	r.POST("/blog/create", func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		fmt.Println(c.Request.PostForm, title, content)
		result := db.Create(&Blog{Title: title, Content: content})
		if errorDB(result, c) {
			return
		}
	})
	r.GET("/blog/get", func(c *gin.Context) {
		var blog Blog
		id, _ := c.GetQuery("id")
		result := db.First(&blog, id)
		if errorDB(result, c) {
			return
		}
		fmt.Println(json.NewEncoder(os.Stdout).Encode(blog))
		c.JSON(http.StatusOK, blog)
	})
	r.GET("/blog/list", func(c *gin.Context) {
		var blogs []Blog
		result := db.Find(&blogs)
		if errorDB(result, c) {
			return
		}
		fmt.Println(json.NewEncoder(os.Stdout).Encode(blogs))
		c.JSON(http.StatusOK, blogs)
	})
}

func main() {
	r := gin.Default()
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&Blog{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	listener(r, db)

	fmt.Println("Database connection and migration successful")
	r.Run(":8080")
}
