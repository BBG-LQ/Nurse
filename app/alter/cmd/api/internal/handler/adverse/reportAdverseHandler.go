// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adverse

import (
	"context"
	"net/http"

	"Nurse/app/alter/cmd/api/internal/logic/adverse"
	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"Nurse/pkg/ctxdata"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 匿名上报
func ReportAdverseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReportAdverseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 关键：handler 层取 IP 塞进 ctx
		ctx := context.WithValue(r.Context(), ctxdata.CtxKeyRemoteIP, httpx.GetRemoteAddr(r))

		// l := adverse.NewReportAdverseLogic(r.Context(), svcCtx)
		l := adverse.NewReportAdverseLogic(ctx, svcCtx)

		resp, err := l.ReportAdverse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
