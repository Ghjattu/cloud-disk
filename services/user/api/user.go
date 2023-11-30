package main

import (
	"flag"
	"fmt"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/config"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/handler"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"

	_ "github.com/joho/godotenv/autoload"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
