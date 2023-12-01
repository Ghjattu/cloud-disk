package oss

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	client      *oss.Client
	bucket      *oss.Bucket
	_bucketName string
	_endpoint   string
)

func Init(bucketName, endpoint, accessKeyID, accessKeySecret string) {
	var err error

	client, err = oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	_bucketName = bucketName
	_endpoint = endpoint
}

func UploadFile(objectKey, localFilePath string) (string, error) {
	err := bucket.PutObjectFromFile(objectKey, localFilePath)
	if err != nil {
		return "", err
	}

	err = bucket.SetObjectACL(objectKey, oss.ACLPublicRead)
	if err != nil {
		bucket.DeleteObject(objectKey)
		return "", err
	}

	url := fmt.Sprintf("https://%s.%s/%s", _bucketName, _endpoint, objectKey)

	return url, nil
}
