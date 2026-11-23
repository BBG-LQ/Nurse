// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adverse

import (
	"net/http"

	"Nurse/app/alter/cmd/api/internal/logic/adverse"
	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 统计图表数据
func AdverseStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StatisticsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := adverse.NewAdverseStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.AdverseStatistics(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
