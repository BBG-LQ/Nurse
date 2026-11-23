// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package safety

import (
	"net/http"

	"Nurse/app/alter/cmd/api/internal/logic/safety"
	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 警报列表
func ListAlertsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListAlertsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := safety.NewListAlertsLogic(r.Context(), svcCtx)
		resp, err := l.ListAlerts(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
