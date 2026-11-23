// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"Nurse/app/usercenter/cmd/api/internal/svc"
	"Nurse/app/usercenter/cmd/api/internal/types"
	"Nurse/app/usercenter/cmd/rpc/usercenter"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line

	//rpc 登录
	result, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &usercenter.LoginReq{

		Password: req.Password,
		Mobile:   req.Mobile,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.LoginResp{
		UserId: resp.UserId,
		Token:  result.Token,
		Role:   result.Role,
	}, nil
}
