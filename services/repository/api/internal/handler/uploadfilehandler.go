package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/logic"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/svc"
	"github.com/Ghjattu/cloud-disk/services/repository/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chunk, chunkHeader, err := r.FormFile("chunk")
		fmt.Println("handler receive chunk, size: ", chunkHeader.Size)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer chunk.Close()

		chunkInfo := r.FormValue("chunk_info")
		var req types.UploadFileReq
		err = json.Unmarshal([]byte(chunkInfo), &req)
		if err != nil {
			fmt.Println("unmarshal chunk info error: ", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		fmt.Printf("handler receive chunk info: %+v\n", req)

		l := logic.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(&req, chunk, chunkHeader.Size)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
