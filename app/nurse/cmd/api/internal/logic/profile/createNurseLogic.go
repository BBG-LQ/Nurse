// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"context"
	"strings"

	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNurseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 录入护士档案
func NewCreateNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNurseLogic {
	return &CreateNurseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNurseLogic) CreateNurse(req *types.CreateNurseReq) (resp *types.CreateNurseResp, err error) {
	// todo: add your logic here and delete this line

	// 1. HTTP 类型 → RPC proto 类型
	rpcReq := &nurse.CreateNurseReq{
		RealName:      strings.TrimSpace(req.RealName),
		Mobile:        strings.TrimSpace(req.Mobile),
		LicenseNo:     strings.TrimSpace(req.LicenseNo),
		LicenseExpire: strings.TrimSpace(req.LicenseExpire),
		SkillTags:     req.SkillTags,
		ServiceArea:   strings.TrimSpace(req.ServiceArea),
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.NurseRpc.CreateNurse(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("调用 RPC CreateNurse 失败: req=%+v, err=%v", req, err)
		return nil, xerr.NewErrMsg(err.Error())
	}

	return &types.CreateNurseResp{
		Id: rpcResp.GetId(),
	}, nil
}
