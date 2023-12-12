# Cloud Disk

## Description
前端使用 React.js，后端使用 go-zero 框架、MySQL 和 Redis 实现的简易网盘系统，最终文件数据保存本地并使用 Nginx 作为静态资源代理，支持文件秒传、分片上传、文件下载功能。

分片的大小默认为 5MB。

## Start
### Nginx
首先确保安装了 Nginx，然后在 `sites-available` 目录下为此项目创建一个配置文件：
```shell
cd /etc/nginx/sites-available
mkdir clouddisk && cd clouddisk
vim clouddisk.conf
```

然后在 `clouddisk.conf` 中输入下列内容：
```
server {
    listen 8080;
    server_name <your-ip>;

    location / {
        root <your-index-path>;
        index index.html;
    }

	location /user {
		proxy_pass http://127.0.0.1:8081;
	}

	location /file {
		# 这一行是限制请求体的大小
	    client_max_body_size 10M;

		proxy_pass http://127.0.0.1:8082;
	}

	# 下载文件时的路由，将 alias 后面的路径改为你的静态文件夹的绝对路径
	location /static/videos {
		alias <your-static-directory-path>;
	}
}
```

接着在 `sites-enabled` 目录下创建一个软链接，并且重新加载配置文件:
```shell
ln -s /etc/nginx/sites-available/clouddisk/clouddisk.conf /etc/nginx/sites-enabled/
/usr/sbin/nginx -s reload
```

### Backend
#### User Service
需要先创建一个 `.env` 文件：
```shell
cd services/user/api
touch .env
```

然后输入下列的环境变量：
```
MYSQL_HOST=
MYSQL_PORT=
MYSQL_USER=
MYSQL_DATABASE=

ACCESS_SECRET=<jwt-secret-key>
ACCESS_EXPIRE=<jwt-expire-seconds>
```

然后执行 `go run user.go` 即可启动用户服务，用户服务运行在 8081 端口。

#### File Service
同样需要先创建一个 `.env` 文件：
```shell
cd services/repository/api
touch .env
```

然后输入下列的环境变量：
```
ACCESS_SECRET=<jwt-secret-key>
ACCESS_EXPIRE=<jwt-expire-seconds>

MYSQL_HOST=
MYSQL_PORT=
MYSQL_USER=
MYSQL_DATABASE=

REDIS_HOST=<ip:port>

# 和 Nginx 配置文件中的路径一致
STATIC_PATH=<your-static-directory-path>
```

然后执行 `go run repository.go` 即可启动文件服务，文件服务运行在 8082 端口。

### Frontend
在 client 文件夹下，创建一个 `.env` 文件：
```
# 5242880=5MB
REACT_APP_CHUNK_SIZE=5242880

# 限制前端同时上传的请求数量
REACT_APP_WINDOW_SIZE=10
```

执行 `npm start` 即可启动前端服务，前端服务运行在 3000 端口。
```shell
cd client
npm start
```

如果想打包前端的话执行 `npm run build` 即可。