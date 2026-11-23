// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package archive

import (
	"context"

	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateElderlyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateElderlyLogic {
	return &UpdateElderlyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateElderlyLogic) UpdateElderly(req *types.UpdateElderlyReq) (resp *types.UpdateElderlyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
