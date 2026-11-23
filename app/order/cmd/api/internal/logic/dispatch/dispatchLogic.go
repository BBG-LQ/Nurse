// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dispatch

import (
	"context"

	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 智能调度
func NewDispatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchLogic {
	return &DispatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchLogic) Dispatch(req *types.DispatchReq) (resp *types.DispatchResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if req.OrderId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}

	// 2. 调用 RPC
	//    派单核心逻辑（状态机、护士选择、负载预占、并发互斥）全部在 RPC 层
	rpcResp, err := l.svcCtx.OrderRpc.Dispatch(l.ctx, &order.DispatchReq{
		OrderId: req.OrderId,
	})
	if err != nil {
		// 不重造错误，交给全局 ErrorHandler
		l.Errorf("调用 Dispatch RPC 失败: orderId=%d, err=%v", req.OrderId, err)
		return nil, err
	}

	// 3. 映射响应
	return &types.DispatchResp{
		NurseId: rpcResp.NurseId,
	}, nil
}
