// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package medication

import (
	"net/http"

	"Nurse/app/elderly/cmd/api/internal/logic/medication"
	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateMedicationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateMedicationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := medication.NewCreateMedicationLogic(r.Context(), svcCtx)
		resp, err := l.CreateMedication(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
