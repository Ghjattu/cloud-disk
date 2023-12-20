// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/login",
				Handler: LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/register",
				Handler: RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/oauth/github/login_url",
				Handler: GetGithubLoginURLHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/oauth/github/callback/:code",
				Handler: GithubCallbackHandler(serverCtx),
			},
		},
	)
}
