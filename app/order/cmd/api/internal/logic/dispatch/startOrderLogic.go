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

type StartOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 开始服务
func NewStartOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartOrderLogic {
	return &StartOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartOrderLogic) StartOrder(req *types.StartReq) (resp *types.StartResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验（path:"id" 由 go-zero 自动解析）
	if req.OrderId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}

	// 2. 前置拦截未登录：RPC 层依赖 ctx uid 做归属校验，
	//    这里先拦住 uid=0 的情况，避免无谓的 RPC 往返
	if ctxdata.GetUidFromCtx(l.ctx) <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.TokenExpireError, "未登录或登录已过期")
	}

	// 3. 调用 RPC
	//    归属校验（uid == order.nurse_id）在 RPC 层完成，
	//    依赖 zrpc 客户端把 jwtUserId 从 ctx 透传到 gRPC metadata
	rpcResp, err := l.svcCtx.OrderRpc.StartOrder(l.ctx, &order.StartOrderReq{
		OrderId: req.OrderId,
	})
	if err != nil {
		l.Errorf("调用 StartOrder RPC 失败: orderId=%d, err=%v", req.OrderId, err)
		return nil, err
	}

	// 4. 映射响应
	return &types.StartResp{Ok: rpcResp.Ok}, nil
}
