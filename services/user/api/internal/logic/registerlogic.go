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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Password: string(hashedPassword),
	}
	err = l.svcCtx.DB.Create(user).Error
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(int64(user.ID), user.Name)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		Token: token,
	}, nil
}
