package logic

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/oss"

	"github.com/Ghjattu/cloud-disk/services/repository/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

var mu sync.Mutex

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
	// resp := &types.UploadFileResp{}
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

	// get the count of saved chunks
	mu.Lock()
	defer mu.Unlock()
	chunkCount, err := l.svcCtx.Redis.Hlen(redisKey)
	if err != nil {
		return nil, errors.New(1, "get chunk count failed")
	}
	if chunkCount == req.TotalChunks {
		// merge chunks
		savedLocalPath := fmt.Sprintf("./%d_%s", currentUserID, req.FileHash)
		defer os.Remove(savedLocalPath)

		err = utils.MergeChunks(l.svcCtx.Redis, redisKey, savedLocalPath, req.FileHash)
		if err != nil {
			return nil, errors.New(1, "merge chunks failed")
		}

		// upload file to oss
		objectKey := fmt.Sprintf("%d_%s%s", currentUserID, req.FileHash, path.Ext(req.FileName))
		ossPath, _ := oss.UploadFile(objectKey, savedLocalPath)

		now := time.Now()
		// save file meta to mysql
		fileModel := &model.File{
			OwnerID:    currentUserID,
			Hash:       req.FileHash,
			Name:       req.FileName,
			Size:       req.FileSize,
			Path:       ossPath,
			UploadTime: now,
		}

		l.svcCtx.DB.Model(&model.File{}).Create(fileModel)

		// delete redis key
		l.svcCtx.Redis.Del(redisKey)

		return &types.UploadFileResp{
			FileSuccess:  true,
			ChunkSuccess: true,
			ChunksCount:  req.TotalChunks,
			FileID:       int64(fileModel.ID),
			FileURL:      ossPath,
			UploadTime:   now.Unix(),
		}, nil
	}

	return &types.UploadFileResp{
		FileSuccess:  false,
		ChunkSuccess: true,
		ChunksCount:  chunkCount,
		FileID:       -1,
		FileURL:      "",
	}, nil
}
