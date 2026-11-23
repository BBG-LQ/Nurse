// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dispatch

import (
	"context"

	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 护士接单
func NewAcceptOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptOrderLogic {
	return &AcceptOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AcceptOrderLogic) AcceptOrder(req *types.AcceptReq) (resp *types.AcceptResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if req.OrderId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}

	// 2. 身份确定：以 JWT uid 为准
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.TokenExpireError, "未登录或登录已过期")
	}
	// req.NurseId 仅做一致性校验，防止前端传错或伪造
	if req.NurseId > 0 && req.NurseId != uid {
		l.Infof("护士ID与登录身份不符: reqNurseId=%d, reqUid=%d", req.NurseId, uid)
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID与登录身份不符")
	}

	// 3. 调用 RPC（用 uid 而非 req.NurseId）
	rpcResp, err := l.svcCtx.OrderRpc.AcceptOrder(l.ctx, &order.AcceptOrderReq{
		OrderId: req.OrderId,
		NurseId: uid,
	})
	if err != nil {
		l.Errorf("调用 AcceptOrder RPC 失败: orderId=%d, nurseId=%d, err=%v",
			req.OrderId, uid, err)
		return nil, err
	}

	// 4. 映射响应
	return &types.AcceptResp{Ok: rpcResp.Ok}, nil
}
