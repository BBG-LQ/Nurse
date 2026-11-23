// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ingest

import (
	"net/http"

	"Nurse/app/vital/cmd/api/internal/logic/ingest"
	"Nurse/app/vital/cmd/api/internal/svc"
	"Nurse/app/vital/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 模拟手环上报
func IngestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IngestReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ingest.NewIngestLogic(r.Context(), svcCtx)
		resp, err := l.Ingest(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
