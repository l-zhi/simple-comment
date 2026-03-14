package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"simple-comment/internal/config"
	"simple-comment/internal/handler"
	"simple-comment/internal/model"
	"simple-comment/internal/repository"
	"simple-comment/internal/service"
)

func New(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong from simple-comment backend",
		})
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		log.Printf("[migrate] AutoMigrate: %v", err)
		panic(err)
	}

	repo := repository.NewCommentRepository(db)
	svc := service.NewCommentService(repo)
	h := handler.NewCommentHandler(svc)

	api := r.Group("/api")
	h.RegisterRoutes(api)

	return r
}
