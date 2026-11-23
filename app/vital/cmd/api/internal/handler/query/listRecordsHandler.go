// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package query

import (
	"net/http"

	"Nurse/app/vital/cmd/api/internal/logic/query"
	"Nurse/app/vital/cmd/api/internal/svc"
	"Nurse/app/vital/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 近 N 天趋势
func ListRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListRecordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := query.NewListRecordsLogic(r.Context(), svcCtx)
		resp, err := l.ListRecords(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
