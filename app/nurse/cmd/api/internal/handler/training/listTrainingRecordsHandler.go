// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package training

import (
	"net/http"

	"Nurse/app/nurse/cmd/api/internal/logic/training"
	"Nurse/app/nurse/cmd/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 个人考核记录
func ListTrainingRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := training.NewListTrainingRecordsLogic(r.Context(), svcCtx)
		resp, err := l.ListTrainingRecords()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
