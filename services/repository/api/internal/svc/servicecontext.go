package svc

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/config"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/workerpool"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *redis.Redis
	StaticPath string // local path for saved files
	JobChan    chan workerpool.Job
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

	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		panic("STATIC_PATH is not set")
	}
	absPath, err := filepath.Abs(staticPath)
	if err != nil {
		panic(err)
	}
	// create the directory if it does not exist
	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// create worker pool
	jobChan := make(chan workerpool.Job, c.WorkerPool.JobChannelSize)
	workerpool.CreateWorkerPool(c.WorkerPool.MaxWorkers, jobChan)

	return &ServiceContext{
		Config:     c,
		DB:         db,
		Redis:      redis.MustNewRedis(c.RedisConf),
		StaticPath: absPath,
		JobChan:    jobChan,
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
