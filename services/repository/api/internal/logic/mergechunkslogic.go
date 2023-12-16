package logic

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type MergeChunksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMergeChunksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MergeChunksLogic {
	return &MergeChunksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MergeChunksLogic) MergeChunks(req *types.MergeChunksReq) (resp *types.MergeChunksResp, err error) {
	currentUserIDStr := fmt.Sprintf("%v", l.ctx.Value("user_id"))
	currentUserID, _ := strconv.ParseInt(currentUserIDStr, 10, 64)

	redisKey := fmt.Sprintf("%d_%s", currentUserID, req.FileHash)
	localFileName := fmt.Sprintf("%d_%s%s", currentUserID, req.FileHash, path.Ext(req.FileName))
	localFilePath := fmt.Sprintf("%s/%s", l.svcCtx.StaticPath, localFileName)
	err = utils.MergeChunks(l.svcCtx.Redis, redisKey, localFilePath, req.FileHash)
	if err != nil {
		return nil, errors.New(1, "merge chunks failed: "+err.Error())
	}

	now := time.Now()
	// save file meta to mysql
	fileModel := &model.File{
		OwnerID:    currentUserID,
		Hash:       req.FileHash,
		Name:       req.FileName,
		Size:       req.FileSize,
		Path:       localFileName,
		UploadTime: now,
	}
	l.svcCtx.DB.Model(&model.File{}).Create(fileModel)

	// delete redis key
	l.svcCtx.Redis.Del(redisKey)

	return &types.MergeChunksResp{
		FileID:     int64(fileModel.ID),
		FileURL:    localFileName,
		UploadTime: now.Unix(),
	}, nil

}
