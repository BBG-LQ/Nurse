// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"net/http"

	"Nurse/app/nurse/cmd/api/internal/logic/profile"
	"Nurse/app/nurse/cmd/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 护士详情
func GetNurseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := profile.NewGetNurseLogic(r.Context(), svcCtx)
		resp, err := l.GetNurse()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
