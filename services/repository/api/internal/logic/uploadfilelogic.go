package logic

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strconv"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/oss"
	"github.com/Ghjattu/cloud-disk/services/repository/utils"

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

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, chunk multipart.File, chunkSize int64) (resp *types.UploadFileResp, err error) {
	fmt.Printf("upload file logic receive chunk num: %d\n", req.ChunkNum)
	currentUserIDStr := fmt.Sprintf("%v", l.ctx.Value("user_id"))
	currentUserID, _ := strconv.ParseInt(currentUserIDStr, 10, 64)

	// consistency check of the chunk
	chunkHash, err := utils.GetMD5Hash(chunk, chunkSize)
	if err != nil {
		resp.StatusMsg = err.Error()
		return
	}
	if chunkHash != req.ChunkHash {
		err = fmt.Errorf("chunk hash mismatch")
		resp.StatusMsg = "chunk hash mismatch"
		return
	}

	// save chunk in redis and set expiration time to 24 hours
	redisKey := fmt.Sprintf("%d_%s", currentUserID, req.FileHash)
	err = utils.SaveChunkInRedis(l.svcCtx.Redis, chunk, redisKey, req.ChunkNum)
	if err != nil {
		resp.StatusMsg = err.Error()
		return
	}

	// get the count of saved chunks
	chunkCount, err := l.svcCtx.Redis.Hlen(redisKey)
	fmt.Printf("chunk count: %d\n\n", chunkCount)
	if err != nil {
		resp.StatusMsg = err.Error()
		return
	}
	if chunkCount == req.TotalChunks {
		// merge chunks
		savedLocalPath := fmt.Sprintf("./%d_%s", currentUserID, req.FileHash)
		defer os.Remove(savedLocalPath)

		err = utils.MergeChunks(l.svcCtx.Redis, redisKey, savedLocalPath, req.FileHash)
		if err != nil {
			resp.StatusMsg = err.Error()
			return
		}

		// upload file to oss
		objectKey := fmt.Sprintf("%d_%s%s", currentUserID, req.FileHash, path.Ext(req.FileName))
		ossPath, _ := oss.UploadFile(objectKey, savedLocalPath)

		// save file meta to mysql
		fileModel := &model.File{
			OwnerID: currentUserID,
			Hash:    req.FileHash,
			Name:    req.FileName,
			Ext:     path.Ext(req.FileName),
			Size:    req.FileSize,
			Path:    ossPath,
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
