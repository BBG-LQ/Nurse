// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package medication

import (
	"net/http"

	"Nurse/app/elderly/cmd/api/internal/logic/medication"
	"Nurse/app/elderly/cmd/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListMedicationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := medication.NewListMedicationsLogic(r.Context(), svcCtx)
		resp, err := l.ListMedications()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
