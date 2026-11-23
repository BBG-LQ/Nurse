package logic

import (
	"context"

	"Nurse/app/usercenter/cmd/rpc/internal/svc"
	"Nurse/app/usercenter/cmd/rpc/usercenter"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListFamilyEldersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFamilyEldersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFamilyEldersLogic {
	return &ListFamilyEldersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFamilyEldersLogic) ListFamilyElders(in *usercenter.ListFamilyEldersReq) (*usercenter.ListFamilyEldersResp, error) {
	// todo: add your logic here and delete this line

	whereBuilder := l.svcCtx.FamilyBindingModel.SelectBuilder().Where(squirrel.Eq{
		"family_user_id": in.FamilyUserId,
	})

	bindings, err := l.svcCtx.FamilyBindingModel.FindAll(l.ctx, whereBuilder, "id ASC")

	if err != nil {
		return nil, err
	}

	// 3. 提取 elderly_id
	elderlyIds := make([]int64, 0, len(bindings))
	for _, b := range bindings {
		elderlyIds = append(elderlyIds, int64(b.ElderlyId))
	}

	return &usercenter.ListFamilyEldersResp{
		ElderlyIds: elderlyIds,
	}, nil
}
