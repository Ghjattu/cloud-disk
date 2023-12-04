package logic

import (
	"context"
	"fmt"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileListLogic) GetFileList() (resp []types.GetFileListResp, err error) {
	currentUserIDStr := fmt.Sprintf("%v", l.ctx.Value("user_id"))

	fileList := make([]model.File, 0)
	err = l.svcCtx.DB.Model(&model.File{}).
		Where("owner_id = ?", currentUserIDStr).
		Find(&fileList).Error
	if err != nil {
		return
	}

	for _, file := range fileList {
		resp = append(resp, types.GetFileListResp{
			FileID:   int64(file.ID),
			FileName: file.Name,
			FileSize: file.Size,
			FileURL:  file.Path,
		})
	}

	return
}
