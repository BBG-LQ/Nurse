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

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get user profile
func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (resp *types.ProfileResp, err error) {
	// todo: add your logic here and delete this line
	//rpc 获取detail
	result, err := l.svcCtx.UsercenterRpc.GetUser(l.ctx, &usercenter.GetUserReq{

		UserId: ctxdata.GetUidFromCtx(l.ctx),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", "req nil")
	}
	return &types.ProfileResp{
		Id:       result.Id,
		RealName: result.RealName,
		Mobile:   result.Mobile,
		Role:     result.Role,
	}, nil
}
