package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sort"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func GetMD5Hash(file multipart.File, fileSize int64) (string, error) {
	file.Seek(0, 0)
	b := make([]byte, fileSize)
	_, err := file.Read(b)
	if err != nil {
		return "", err
	}

	fileHash := fmt.Sprintf("%x", md5.Sum(b))
	return fileHash, nil
}

func SaveUploadedFile(file multipart.File, localPath string) error {
	// create a local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// copy file to local
	file.Seek(0, 0)
	_, err = io.Copy(localFile, file)
	if err != nil {
		return err
	}

	return nil
}

func SaveChunkInRedis(redis *redis.Redis, chunk multipart.File, key string, field int) error {
	chunk.Seek(0, 0)
	bytes, err := io.ReadAll(chunk)
	if err != nil {
		return err
	}

	fieldStr := fmt.Sprintf("%d", field)
	err = redis.Hset(key, fieldStr, string(bytes))
	if err != nil {
		return err
	}

	redis.Expire(key, 60*60*24)

	return nil
}

func MergeChunks(redis *redis.Redis, key, savedLocalPath, fileHash string, totalChunks int64) error {
	fields, err := redis.Hkeys(key)
	if err != nil {
		return err
	}

	if int64(len(fields)) != totalChunks {
		return fmt.Errorf("incomplete chunks")
	}

	// sort by chunk number
	sort.Slice(fields, func(i, j int) bool {
		numA, _ := strconv.Atoi(fields[i])
		numB, _ := strconv.Atoi(fields[j])
		return numA < numB
	})

	// create a local file
	localFile, err := os.OpenFile(savedLocalPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// merge chunks
	for _, field := range fields {
		chunk, err := redis.Hget(key, field)
		if err != nil {
			return err
		}

		_, err = localFile.Write([]byte(chunk))
		if err != nil {
			return err
		}
	}

	// consistency check of the merged file
	localFileInfo, _ := localFile.Stat()
	localFileSize := localFileInfo.Size()
	localFileHash, _ := GetMD5Hash(localFile, localFileSize)
	if localFileHash != fileHash {
		return fmt.Errorf("file consistency check failed")
	}

	return nil
}
