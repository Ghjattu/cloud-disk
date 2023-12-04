package handler

import (
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckFileExistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckFileExistReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCheckFileExistLogic(r.Context(), svcCtx)
		resp, err := l.CheckFileExist(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
