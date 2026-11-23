// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"Nurse/app/usercenter/cmd/api/internal/svc"
	"Nurse/app/usercenter/cmd/api/internal/types"
	"Nurse/app/usercenter/cmd/rpc/usercenter"
	"Nurse/pkg/ctxdata"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update user profile
func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) (resp *types.UpdateProfileResp, err error) {
	// todo: add your logic here and delete this line

	result, err := l.svcCtx.UsercenterRpc.UpdateProfile(l.ctx, &usercenter.UpdateProfileReq{

		UserId:   ctxdata.GetUidFromCtx(l.ctx),
		RealName: req.RealName,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.UpdateProfileResp{
		Ok: result.Ok,
	}, nil
}
