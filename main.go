package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
	r.GET("/hoge", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "fuga",
		})
	})
	fmt.Println("Database connection and migration successful")
	r.Run(":8080")
}
