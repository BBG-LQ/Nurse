// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package safety

import (
	"context"

	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SosNurseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 护士一键报警
func NewSosNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SosNurseLogic {
	return &SosNurseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SosNurseLogic) SosNurse(req *types.SosNurseReq) (resp *types.SosNurseResp, err error) {
	// todo: add your logic here and delete this line

	// 1. API 层基础校验（快速失败）
	if req.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不能为空")
	}
	if req.OrderId < 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}

	// 2. types.Req -> rpc Req 转换
	rpcResp, err := l.svcCtx.AlertRpc.SosNurse(l.ctx, &alert.SosNurseReq{
		NurseId: req.NurseId,
		OrderId: req.OrderId,
		Content: req.Content,
	})
	if err != nil {
		l.Errorf("call rpc SosNurse failed, nurse_id=%d, order_id=%d, err=%v", req.NurseId, req.OrderId, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "sos nurse failed: %v", err)
	}

	// 3. rpc Resp -> types.Resp 转换
	return &types.SosNurseResp{
		AlertId: rpcResp.AlertId,
	}, nil
}
