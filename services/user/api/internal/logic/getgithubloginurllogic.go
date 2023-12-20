package logic

import (
	"context"
	"fmt"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGithubLoginURLLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGithubLoginURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGithubLoginURLLogic {
	return &GetGithubLoginURLLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGithubLoginURLLogic) GetGithubLoginURL() (resp *types.GetGithubLoginURLResp, err error) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user",
		l.svcCtx.Config.OAuthGithub.ClientID,
		l.svcCtx.Config.OAuthGithub.RedirectURL,
	)

	return &types.GetGithubLoginURLResp{
		URL: url,
	}, nil
}
