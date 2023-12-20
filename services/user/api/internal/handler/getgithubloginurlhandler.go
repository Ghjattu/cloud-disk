package handler

import (
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/user/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/user/api/internal/svc"
	xhttp "github.com/zeromicro/x/http"
)

func GetGithubLoginURLHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetGithubLoginURLLogic(r.Context(), svcCtx)
		resp, err := l.GetGithubLoginURL()
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
