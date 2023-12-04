# Cloud Disk

## Description
前端使用 React.js，后端使用 go-zero 框架、MySQL 和 Redis 实现的简易网盘系统，最终文件数据保存在阿里云 OSS 中，支持文件秒传、分片上传、文件下载功能。

## Start
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
OSS_BUCKET_NAME=
OSS_ENDPOINT=
OSS_ACCESS_KEY_ID=
OSS_ACCESS_KEY_SECRET=

ACCESS_SECRET=<jwt-secret-key>
ACCESS_EXPIRE=<jwt-expire-seconds>

MYSQL_HOST=
MYSQL_PORT=
MYSQL_USER=
MYSQL_DATABASE=

REDIS_HOST=<ip:port>
```

然后执行 `go run repository.go` 即可启动文件服务，文件服务运行在 8082 端口。

### Frontend
在 client 文件夹下执行 `npm start` 即可启动前端服务。
```shell
cd client
npm start
```
