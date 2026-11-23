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

type BindFamilyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// bind
func NewBindFamilyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindFamilyLogic {
	return &BindFamilyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindFamilyLogic) BindFamily(req *types.BindFamilyReq) (resp *types.BindFamilyResp, err error) {
	// todo: add your logic here and delete this line

	//rpc 获取detail
	result, err := l.svcCtx.UsercenterRpc.BindFamily(l.ctx, &usercenter.BindFamilyReq{

		FamilyUserId: ctxdata.GetUidFromCtx(l.ctx),
		ElderlyId:    req.ElderlyId,
		Relation:     req.Relation,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.BindFamilyResp{
		Ok: result.Ok,
	}, nil
}
