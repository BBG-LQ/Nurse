// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package medication

import (
	"context"

	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMedicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMedicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMedicationLogic {
	return &CreateMedicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMedicationLogic) CreateMedication(req *types.CreateMedicationReq) (resp *types.CreateMedicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
