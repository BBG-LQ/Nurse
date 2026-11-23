package logic

import (
	"context"

	"Nurse/app/usercenter/cmd/rpc/internal/svc"
	"Nurse/app/usercenter/cmd/rpc/usercenter"
	"Nurse/app/usercenter/model"
	"Nurse/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *usercenter.GetUserReq) (*usercenter.GetUserResp, error) {
	// todo: add your logic here and delete this line
	//1.查询是否存在 根据id
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, uint64(in.UserId))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GetUserInfo find user db err , id:%d , err:%v", in.UserId, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, " id:%d ", in.UserId)
	}

	// 2.返回结构体
	var respUser usercenter.GetUserResp

	_ = copier.Copy(&respUser, user)

	return &usercenter.GetUserResp{

		Id:       respUser.Id,
		Mobile:   respUser.Mobile,
		RealName: respUser.RealName,

		Role:   respUser.Role,
		Status: respUser.Status,
	}, nil
}
