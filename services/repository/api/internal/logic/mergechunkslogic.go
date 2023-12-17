package logic

import (
	"context"
	"fmt"
	"math"
	"path"
	"strconv"
	"time"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"
	"github.com/Ghjattu/cloud-disk/services/repository/workerpool"

	"github.com/zeromicro/go-zero/core/logx"
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

	prefix := fmt.Sprintf("%d_%s", currentUserID, req.FileHash)
	lockName := fmt.Sprintf("%s_%s", prefix, "merge_lock")
	localFileName := fmt.Sprintf("%d_%s%s", currentUserID, req.FileHash, path.Ext(req.FileName))
	localFilePath := fmt.Sprintf("%s/%s", l.svcCtx.StaticPath, localFileName)

	success, err := l.svcCtx.Redis.Setnx(lockName, "true")
	if err != nil {
		return nil, err
	}
	// only the first request can get the lock
	if success {
		// one hour
		l.svcCtx.Redis.Expire(lockName, 60*60)

		mergeJob := workerpool.Job{
			Redis:          l.svcCtx.Redis,
			RedisKeyPrefix: prefix,
			TotalChunks:    req.TotalChunks,
			LocalFilePath:  localFilePath,
			FileHash:       req.FileHash,
		}
		l.svcCtx.JobChan <- mergeJob
	}

	checkNum := int(math.Ceil(float64(l.svcCtx.Config.Timeout)/1000)) - 1
	redisKey := fmt.Sprintf("%s_chunks", prefix)
	resultKey := fmt.Sprintf("%s_result", prefix)
	for i := 0; i < checkNum; i++ {
		// check merge result
		value, err := l.svcCtx.Redis.Get(resultKey)
		if err == nil && value != "" {
			if value == "success" {
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

				// delete keys
				l.svcCtx.Redis.Del(redisKey, resultKey, lockName)

				return &types.MergeChunksResp{
					FileID:     int64(fileModel.ID),
					FileURL:    localFileName,
					UploadTime: now.Unix(),
				}, nil
			} else {
				return nil, fmt.Errorf("merge failed: %s", value)
			}
		}

		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("timeout")
}
