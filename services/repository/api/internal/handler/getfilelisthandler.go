package handler

import (
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	xhttp "github.com/zeromicro/x/http"
)

func GetFileListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetFileListLogic(r.Context(), svcCtx)
		resp, err := l.GetFileList()
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
