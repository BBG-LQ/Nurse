// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package archive

import (
	"context"

	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListElderlyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListElderlyLogic {
	return &ListElderlyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListElderlyLogic) ListElderly(req *types.ListElderlyReq) (resp *types.ListElderlyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
