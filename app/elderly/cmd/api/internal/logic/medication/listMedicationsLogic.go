// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package medication

import (
	"context"

	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMedicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMedicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMedicationsLogic {
	return &ListMedicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMedicationsLogic) ListMedications() (resp *types.ListMedicationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
