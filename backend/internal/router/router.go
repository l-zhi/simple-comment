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

	if err := db.AutoMigrate(&model.Comment{}, &model.Article{}); err != nil {
		log.Printf("[migrate] AutoMigrate: %v", err)
		panic(err)
	}

	commentRepo := repository.NewCommentRepository(db)
	commentSvc := service.NewCommentService(commentRepo)
	commentH := handler.NewCommentHandler(commentSvc)

	articleRepo := repository.NewArticleRepository(db)
	articleSvc := service.NewArticleService(articleRepo)
	articleH := handler.NewArticleHandler(articleSvc)

	api := r.Group("/api")
	commentH.RegisterRoutes(api)
	articleH.RegisterRoutes(api)

	return r
}
