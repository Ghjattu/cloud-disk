package handler

import (
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
)

func GithubCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GithubCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGithubCallbackLogic(r.Context(), svcCtx)
		resp, err := l.GithubCallback(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
