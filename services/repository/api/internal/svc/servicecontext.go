package svc

import (
	"fmt"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/config"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(getDSN(c)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.File{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  redis.MustNewRedis(c.RedisConf),
	}
}

func getDSN(c config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQL.User,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Database,
	)
}
