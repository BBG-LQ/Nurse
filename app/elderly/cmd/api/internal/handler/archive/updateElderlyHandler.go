// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package archive

import (
	"net/http"

	"Nurse/app/elderly/cmd/api/internal/logic/archive"
	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateElderlyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateElderlyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := archive.NewUpdateElderlyLogic(r.Context(), svcCtx)
		resp, err := l.UpdateElderly(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
