package logic

import (
	"context"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	user := &model.User{}
	err = l.svcCtx.DB.Model(&model.User{}).Where("id = ?", req.ID).First(user).Error
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResp{
		Name: user.Name,
	}, nil
}
