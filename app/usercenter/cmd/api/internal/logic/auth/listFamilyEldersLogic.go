// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"Nurse/app/usercenter/cmd/api/internal/svc"
	"Nurse/app/usercenter/cmd/api/internal/types"
	"Nurse/app/usercenter/cmd/rpc/usercenter"
	"Nurse/pkg/ctxdata"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListFamilyEldersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// listFamilyElders
func NewListFamilyEldersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFamilyEldersLogic {
	return &ListFamilyEldersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFamilyEldersLogic) ListFamilyElders() (resp *types.FamilyEldersResp, err error) {
	// todo: add your logic here and delete this line

	//rpc 获取
	result, err := l.svcCtx.UsercenterRpc.ListFamilyElders(l.ctx, &usercenter.ListFamilyEldersReq{
		FamilyUserId: ctxdata.GetUidFromCtx(l.ctx),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", "req nil")
	}

	return &types.FamilyEldersResp{
		ElderlyIds: result.ElderlyIds,
	}, nil
}
