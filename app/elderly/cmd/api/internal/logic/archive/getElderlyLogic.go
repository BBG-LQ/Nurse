// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package archive

import (
	"context"

	"Nurse/app/elderly/cmd/api/internal/svc"
	"Nurse/app/elderly/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetElderlyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetElderlyLogic {
	return &GetElderlyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetElderlyLogic) GetElderly() (resp *types.GetElderlyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
