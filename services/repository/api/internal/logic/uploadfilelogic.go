package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/oss"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (resp *types.UploadFileResp, err error) {
	currentUserIDStr := fmt.Sprintf("%v", l.ctx.Value("user_id"))

	// get md5 hash of file
	b := make([]byte, fileHeader.Size)
	_, err = file.Read(b)
	if err != nil {
		return nil, err
	}
	fileHash := fmt.Sprintf("%x", md5.Sum(b))

	// check if file exists
	var count int64 = 0
	err = l.svcCtx.DB.Model(&model.File{}).
		Where("id = ? AND hash = ?", currentUserIDStr, fileHash).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		// file exists, return
		return nil, fmt.Errorf("file already exists")
	}

	// Save video to local.
	publishTimeStr := time.Now().Format("2006-01-02-15:04:05")
	finalFileName := fmt.Sprintf("%s_%s_%s", currentUserIDStr, publishTimeStr, fileHeader.Filename)
	fileSavedLocalPath := filepath.Join("./", finalFileName)

	// create a local file
	localFile, err := os.Create(fileSavedLocalPath)
	if err != nil {
		fmt.Println("os.Create err: ", err)
		return nil, err
	}
	// defer os.Remove(fileSavedLocalPath)
	defer localFile.Close()

	// copy file to local
	file.Seek(0, 0)
	_, err = io.Copy(localFile, file)
	if err != nil {
		return nil, err
	}

	// Upload video to OSS.
	if err := oss.UploadFile(finalFileName, fileSavedLocalPath); err != nil {
		return nil, err
	}

	// save file info to mysql
	bucketName := l.svcCtx.Config.OSS.BucketName
	endpoint := l.svcCtx.Config.OSS.Endpoint
	ossPath := fmt.Sprintf("https://%s.%s/%s", bucketName, endpoint, finalFileName)

	fileModel := &model.File{
		Hash: fileHash,
		Name: fileHeader.Filename,
		Ext:  path.Ext(fileHeader.Filename),
		Size: fileHeader.Size,
		Path: ossPath,
	}
	err = l.svcCtx.DB.Model(&model.File{}).Create(fileModel).Error
	if err != nil {
		return nil, err
	}

	return &types.UploadFileResp{
		FileID: int64(fileModel.ID),
	}, nil
}
