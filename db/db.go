package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open() (*gorm.DB, *sql.DB, *redis.Client, error) {
	sqlDB := PostgresqlOpen()
	var err error
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, nil, err
	}
	redisClient := RedisNewClient()
	return gormDB, sqlDB, redisClient, nil
}
