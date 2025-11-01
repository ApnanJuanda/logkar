package config

import (
	"bsnack/app/controller/root"
	"bsnack/config/collection"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Router(DB *gorm.DB, redisClient *redis.Client) error {
	router := gin.Default()
	corsConfig(router)

	router.GET("/", root.Index)

	api := router.Group("/api")
	collection.ApiRouter(DB, redisClient, api)

	if err := router.Run(); err != nil {
		return err
	}
	return nil
}
