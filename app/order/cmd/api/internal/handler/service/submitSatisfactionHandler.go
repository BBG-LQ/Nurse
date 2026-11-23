// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package service

import (
	"net/http"

	"Nurse/app/order/cmd/api/internal/logic/service"
	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 满意度评分
func SubmitSatisfactionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SatisfactionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := service.NewSubmitSatisfactionLogic(r.Context(), svcCtx)
		resp, err := l.SubmitSatisfaction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
