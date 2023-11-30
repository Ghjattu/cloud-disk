package logic

import (
	"context"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/user/model"
	"github.com/Ghjattu/cloud-disk/services/user/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	user := &model.User{}
	err = l.svcCtx.DB.Model(&model.User{}).Where("name = ?", req.Name).First(user).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	accessSecret := l.svcCtx.Config.Auth.AccessSecret
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	token, err := utils.GenerateToken(accessSecret, accessExpire, int64(user.ID), user.Name)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token: token,
	}, nil
}
