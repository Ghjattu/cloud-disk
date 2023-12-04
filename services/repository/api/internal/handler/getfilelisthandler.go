package handler

import (
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFileListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetFileListLogic(r.Context(), svcCtx)
		resp, err := l.GetFileList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
