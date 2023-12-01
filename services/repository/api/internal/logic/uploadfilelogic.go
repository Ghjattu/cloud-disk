package logic

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/oss"
	"github.com/Ghjattu/cloud-disk/services/repository/utils"
	"gorm.io/gorm"

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
	currentUserID, _ := strconv.Atoi(currentUserIDStr)

	// get md5 hash of file
	fileHash, err := utils.GetMD5Hash(file, fileHeader)
	if err != nil {
		return nil, err
	}

	// check if file exists
	existedFile := &model.File{}
	err = l.svcCtx.DB.Model(&model.File{}).
		Where("owner_id = ? AND hash = ?", currentUserID, fileHash).
		First(&existedFile).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == nil {
		// file exists
		return &types.UploadFileResp{
			FileID:  int64(existedFile.ID),
			FileURL: existedFile.Path,
		}, nil
	}

	// Save video to local.
	finalFileName := fmt.Sprintf("%s_%s", fileHash, fileHeader.Filename)
	fileSavedLocalPath := filepath.Join("./", finalFileName)
	err = utils.SaveUploadedFile(file, fileSavedLocalPath)
	defer os.Remove(fileSavedLocalPath)
	if err != nil {
		return nil, err
	}

	// Upload video to OSS.
	ossPath, err := oss.UploadFile(finalFileName, fileSavedLocalPath)
	if err != nil {
		return nil, err
	}

	// save file info to mysql
	fileModel := &model.File{
		OwnerID: int64(currentUserID),
		Hash:    fileHash,
		Name:    fileHeader.Filename,
		Ext:     path.Ext(fileHeader.Filename),
		Size:    fileHeader.Size,
		Path:    ossPath,
	}
	err = l.svcCtx.DB.Model(&model.File{}).Create(fileModel).Error
	if err != nil {
		return nil, err
	}

	return &types.UploadFileResp{
		FileID:  int64(fileModel.ID),
		FileURL: fileModel.Path,
	}, nil
}
