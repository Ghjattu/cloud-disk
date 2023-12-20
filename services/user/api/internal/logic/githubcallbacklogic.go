package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/types"
	"github.com/Ghjattu/cloud-disk/services/user/model"
	"github.com/Ghjattu/cloud-disk/services/user/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GithubCallback struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUserInfo struct {
	UserName string `json:"login"`
}

type GithubCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGithubCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GithubCallbackLogic {
	return &GithubCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GithubCallbackLogic) GithubCallback(req *types.GithubCallbackReq) (resp *types.GithubCallbackResp, err error) {
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		l.svcCtx.Config.OAuthGithub.ClientID,
		l.svcCtx.Config.OAuthGithub.ClientSecret,
		req.Code,
	)

	request, _ := http.NewRequest("POST", url, bytes.NewReader([]byte{}))
	request.Header.Set("Accept", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBody, _ := io.ReadAll(response.Body)
	var githubCallback GithubCallback
	err = json.Unmarshal(respBody, &githubCallback)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Printf("%+v\n", githubCallback)

	// get user info
	url = "https://api.github.com/user"
	request, _ = http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	request.Header.Set("Authorization", "Bearer "+githubCallback.AccessToken)

	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	respBody, _ = io.ReadAll(response.Body)
	var githubUserInfo GithubUserInfo
	fmt.Println(string(respBody))
	err = json.Unmarshal(respBody, &githubUserInfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Printf("%+v\n", githubUserInfo)

	// save user in database
	user := &model.User{
		Name:        githubUserInfo.UserName,
		AccessToken: githubCallback.AccessToken,
	}
	l.svcCtx.DB.Create(user)

	// generate jwt
	accessSecret := l.svcCtx.Config.Auth.AccessSecret
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	token, err := utils.GenerateToken(accessSecret, accessExpire, int64(user.ID), user.Name)
	if err != nil {
		return nil, err
	}

	return &types.GithubCallbackResp{
		UserName: githubUserInfo.UserName,
		Token:    token,
	}, nil
}
