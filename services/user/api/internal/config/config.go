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
}
