package logic

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"

	"github.com/Ghjattu/cloud-disk/services/repository/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
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

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, chunk multipart.File, chunkSize int64) (*types.UploadFileResp, error) {
	currentUserIDStr := fmt.Sprintf("%v", l.ctx.Value("user_id"))
	currentUserID, _ := strconv.ParseInt(currentUserIDStr, 10, 64)

	// consistency check of the chunk
	chunkHash, err := utils.GetMD5Hash(chunk, chunkSize)
	if err != nil {
		return nil, errors.New(1, "calculate chunk hash failed: "+err.Error())
	}
	if chunkHash != req.ChunkHash {
		return nil, errors.New(1, "chunk consistency check failed")
	}

	// save chunk in redis and set expiration time to 24 hours
	redisKey := fmt.Sprintf("%d_%s", currentUserID, req.FileHash)
	err = utils.SaveChunkInRedis(l.svcCtx.Redis, chunk, redisKey, req.ChunkNum)
	if err != nil {
		return nil, errors.New(1, "save chunk in redis failed")
	}

	return &types.UploadFileResp{
		ChunkSuccess: true,
	}, nil
}
