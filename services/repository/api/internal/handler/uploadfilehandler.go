package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	xhttp "github.com/zeromicro/x/http"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chunk, chunkHeader, err := r.FormFile("chunk")
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}
		defer chunk.Close()

		chunkInfo := r.FormValue("chunk_info")
		var req types.UploadFileReq
		err = json.Unmarshal([]byte(chunkInfo), &req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(&req, chunk, chunkHeader.Size)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
