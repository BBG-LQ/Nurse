package logic

import (
	"context"
	"strings"

	"Nurse/app/usercenter/cmd/rpc/internal/svc"
	"Nurse/app/usercenter/cmd/rpc/usercenter"
	"Nurse/app/usercenter/model"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProfileLogic) UpdateProfile(in *usercenter.UpdateProfileReq) (*usercenter.UpdateProfileResp, error) {
	// todo: add your logic here and delete this line

	//1.查询是否存在 根据id
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, uint64(in.UserId))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GetUserInfo find user db err , id:%d , err:%v", in.UserId, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, " id:%d ", in.UserId)
	}

	// 2. 参数校验（只允许非空且长度合理的 real_name）
	realName := strings.TrimSpace(in.RealName)
	if realName == "" {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbUpdateAffectedZeroError), "real_name cannot be empty")
	}
	if n := len([]rune(realName)); n > 32 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RequestParamError), "real_name too long, len:%d", n)
	}

	// 无变更，直接返回成功（避免无意义写库 / 触发 version 冲突）
	if user.RealName == realName {
		return &usercenter.UpdateProfileResp{Ok: true}, nil
	}

	// 3. 更新（go-zero 生成的 Update 使用乐观锁：where id=? and version=?）
	user.RealName = realName
	user.Version++ // 新的 version 传给 Update，Update 内部会把旧 version 放进 where

	result, err := l.svcCtx.SysUserModel.Update(l.ctx, nil, user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError),
			"UpdateProfile update user db err, id:%d, err:%v", in.UserId, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError),
			"UpdateProfile rows affected err, id:%d, err:%v", in.UserId, err)
	}
	if rows == 0 {
		// 乐观锁失败：查询后、更新前被别人改过，可提示重试
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError),
			"UpdateProfile optimistic lock conflict, id:%d, version:%d", in.UserId, user.Version)
	}

	return &usercenter.UpdateProfileResp{
		Ok: true,
	}, nil
}
