package main

import (
	"fmt"
	models "gin-todo-app/src/domain/models"
	database "gin-todo-app/src/infra/database"
	repository "gin-todo-app/src/infra/database/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	db, err := database.ConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// リポジトリの初期化
	blogRepo := repository.NewBlogRepository(db)

	//　マイグレーション
	err = db.AutoMigrate(&models.Blog{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	engine.Static("/static", "./Static")

	// htmlのファイルを全て読み込む
	engine.LoadHTMLGlob("src/infra/http/public/*")

	engine.GET("/index", func(c *gin.Context) {
		var blogs []*models.Blog
		blogs, err := blogRepo.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get blogs"})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "HOME",
			"blogs": blogs,
		})
	})

	engine.GET("/edit", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get blog"})
			return
		}
		var blog *models.Blog
		blog, _ = blogRepo.GetByID(c.Request.Context(), uint(id))
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"title": "Edit",
			"blog":  blog,
		})
	})

	engine.GET("/blog/destroy", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get blog"})
			return
		}
		blogRepo.Delete(c.Request.Context(), uint(id))
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	engine.POST("/blog/update", func(c *gin.Context) {
		id, err := strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get blog"})
			return
		}
		var blog *models.Blog
		blog, _ = blogRepo.GetByID(c.Request.Context(), uint(id))
		blog.Title = c.PostForm("title")
		blog.Content = c.PostForm("content")
		blogRepo.Update(c.Request.Context(), blog)
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	engine.POST("/blog/create", func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		blogRepo.Create(c, &models.Blog{Title: title, Content: content})
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	fmt.Println("Database connection and migration successful")
	engine.Run(":8080")
}
