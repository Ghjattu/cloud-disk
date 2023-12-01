package logic

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/repository/model"

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
	if err != nil {
		return &types.CheckFileExistResp{
			Exist: false,
		}, nil
	}

	return &types.CheckFileExistResp{
		Exist:   true,
		FileID:  int64(fileModel.ID),
		FileURL: fileModel.Path,
	}, nil
}
