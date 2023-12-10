package oss

import (
	"fmt"
	"io"
	"os"

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

// chunk size: 5MB
func UploadFileInChunks(objectKey, localFilePath string) (string, error) {
	chunks, err := oss.SplitFileByPartSize(localFilePath, 5*1024*1024)
	if err != nil {
		return "", err
	}
	fd, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	imur, err := bucket.InitiateMultipartUpload(objectKey)
	if err != nil {
		return "", err
	}
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, io.SeekStart)
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			return "", err
		}
		parts = append(parts, part)
	}

	objectACL := oss.ObjectACL(oss.ACLPublicRead)

	_, err = bucket.CompleteMultipartUpload(imur, parts, objectACL)
	if err != nil {
		return "", err
	}

	// fmt.Println("cmur:", cmur.Location)

	url := fmt.Sprintf("https://%s.%s/%s", _bucketName, _endpoint, objectKey)

	return url, nil
}
