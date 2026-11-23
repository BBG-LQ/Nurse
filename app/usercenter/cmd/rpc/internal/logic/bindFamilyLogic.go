package logic

import (
	"context"
	"database/sql"

	"Nurse/app/usercenter/cmd/rpc/internal/svc"
	"Nurse/app/usercenter/cmd/rpc/usercenter"
	"Nurse/app/usercenter/model"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type BindFamilyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindFamilyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindFamilyLogic {
	return &BindFamilyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindFamilyLogic) BindFamily(in *usercenter.BindFamilyReq) (*usercenter.BindFamilyResp, error) {
	// todo: add your logic here and delete this line
	var ok bool
	//1.查询是否存在 根据id
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, uint64(in.FamilyUserId))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GetUserInfo find user db err , id:%d , err:%v", in.FamilyUserId, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, " id:%d ", in.FamilyUserId)
	}

	if err := l.svcCtx.FamilyBindingModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {

		//2.insert进入user
		insertresult, err := l.svcCtx.FamilyBindingModel.Insert(l.ctx, session, &model.FamilyBinding{
			FamilyUserId: uint64(in.FamilyUserId),
			ElderlyId:    uint64(in.ElderlyId),
			Relation: sql.NullString{
				String: in.Relation,
				Valid:  in.Relation != "", // 空串视为 NULL（按业务需要决定）
			},
		})

		if insertresult != nil {
			ok = true
		}

		if err != nil {

		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &usercenter.BindFamilyResp{
		Ok: ok,
	}, nil
}
