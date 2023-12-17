package workerpool

import (
	"fmt"

	"github.com/Ghjattu/cloud-disk/services/repository/utils"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Job struct {
	Redis          *redis.Redis
	RedisKeyPrefix string
	TotalChunks    int64
	LocalFilePath  string
	FileHash       string
}

func worker(jobChan <-chan Job) {
	for job := range jobChan {
		redisKey := fmt.Sprintf("%s_chunks", job.RedisKeyPrefix)
		err := utils.MergeChunks(job.Redis, redisKey, job.LocalFilePath, job.FileHash, job.TotalChunks)
		resultKey := fmt.Sprintf("%s_result", job.RedisKeyPrefix)
		if err != nil {
			job.Redis.Set(resultKey, err.Error())
		} else {
			job.Redis.Set(resultKey, "success")
		}
	}
}

func CreateWorkerPool(maxWorkers int64, jobChan <-chan Job) {
	for i := 0; i < int(maxWorkers); i++ {
		go worker(jobChan)
	}
}
