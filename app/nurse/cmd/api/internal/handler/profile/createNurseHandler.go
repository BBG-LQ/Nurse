// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"net/http"

	"Nurse/app/nurse/cmd/api/internal/logic/profile"
	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 录入护士档案
func CreateNurseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateNurseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := profile.NewCreateNurseLogic(r.Context(), svcCtx)
		resp, err := l.CreateNurse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
