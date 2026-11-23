// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"net/http"

	"Nurse/app/usercenter/cmd/api/internal/logic/auth"
	"Nurse/app/usercenter/cmd/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// listFamilyElders
func ListFamilyEldersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewListFamilyEldersLogic(r.Context(), svcCtx)
		resp, err := l.ListFamilyElders()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
