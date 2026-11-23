// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package training

import (
	"net/http"

	"Nurse/app/nurse/cmd/api/internal/logic/training"
	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 三次不合格名单
func ListUnqualifiedNursesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnqualifiedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := training.NewListUnqualifiedNursesLogic(r.Context(), svcCtx)
		resp, err := l.ListUnqualifiedNurses(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
