package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	client *oss.Client
	bucket *oss.Bucket
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
}

// UploadFile upload file to oss and set public read permission.
//
//	@param objectKey string
//	@param localFilePath string
//	@return error
func UploadFile(objectKey, localFilePath string) error {
	err := bucket.PutObjectFromFile(objectKey, localFilePath)
	if err != nil {
		return err
	}

	err = bucket.SetObjectACL(objectKey, oss.ACLPublicRead)
	if err != nil {
		bucket.DeleteObject(objectKey)
		return err
	}

	return nil
}
