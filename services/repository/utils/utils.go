package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func GetMD5Hash(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	b := make([]byte, fileHeader.Size)
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
