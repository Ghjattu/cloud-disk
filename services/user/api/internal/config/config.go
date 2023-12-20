package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySQL struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	Auth struct { // JWT 认证需要的密钥和过期时间配置
		AccessSecret string
		AccessExpire int64
	}
	OAuthGithub struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
}
