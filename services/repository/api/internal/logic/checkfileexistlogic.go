package logic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckFileExistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckFileExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFileExistLogic {
	return &CheckFileExistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckFileExistLogic) CheckFileExist(req *types.CheckFileExistReq) (resp *types.CheckFileExistResp, err error) {
	currentUserIDStr := fmt.Sprintf("%v", l.ctx.Value("user_id"))
	currentUserID, _ := strconv.Atoi(currentUserIDStr)

	// check file exist
	fileModel := &model.File{}
	err = l.svcCtx.DB.Model(&model.File{}).
		Where("owner_id = ? AND hash = ?", currentUserID, req.Hash).
		First(fileModel).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// server error
		return nil, err
	}
	if err == nil {
		// file already exist
		return &types.CheckFileExistResp{
			Exist:      true,
			FileID:     int64(fileModel.ID),
			FileURL:    fileModel.Path,
			ChunksHash: []string{},
		}, nil
	}

	// now, file not exist
	// retrieve chunks hash from redis
	redisKey := fmt.Sprintf("%d_%s", currentUserID, req.Hash)
	chunksHash, err := l.svcCtx.Redis.Hkeys(redisKey)
	if err != nil {
		return nil, err
	}

	return &types.CheckFileExistResp{
		Exist:      false,
		FileID:     -1,
		FileURL:    "",
		ChunksHash: chunksHash,
	}, nil
}
