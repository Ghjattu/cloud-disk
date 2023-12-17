package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct { // JWT 认证需要的密钥和过期时间配置
		AccessSecret string
		AccessExpire int64
	}
	OSS struct {
		BucketName      string
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
	}
	MySQL struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	RedisConf  redis.RedisConf
	WorkerPool struct {
		MaxWorkers     int64
		JobChannelSize int64
	}
}
