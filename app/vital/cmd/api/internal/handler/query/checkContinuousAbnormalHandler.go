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

// 连续异常判定
func CheckContinuousAbnormalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckAbnormalReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := query.NewCheckContinuousAbnormalLogic(r.Context(), svcCtx)
		resp, err := l.CheckContinuousAbnormal(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
