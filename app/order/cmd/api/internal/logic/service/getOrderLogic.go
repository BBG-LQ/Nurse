// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package service

import (
	"context"

	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 订单详情
func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (resp *types.OrderDetail, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验（防御性：go-zero 已解析 path，但仍拦截非正数）
	if req.Id <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}

	// 2. 调用 RPC 查询订单
	rpcResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderReq{
		OrderId: req.Id,
	})
	if err != nil {
		// 不重造错误，交给全局 ErrorHandler 统一响应
		l.Errorf("调用 GetOrder RPC 失败: orderId=%d, err=%v", req.Id, err)
		return nil, err
	}
	if rpcResp.Order == nil {
		return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "订单不存在")
	}

	// 3. 越权校验（可选）：老人本人或订单护士才能查看；uid=0（管理端）放行
	if uid := ctxdata.GetUidFromCtx(l.ctx); uid > 0 {
		o := rpcResp.Order
		if o.ElderlyId != uid && o.NurseId != uid {
			l.Infof("越权查看订单: orderId=%d, elderlyId=%d, nurseId=%d, reqUid=%d",
				req.Id, o.ElderlyId, o.NurseId, uid)
			return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权查看该订单")
		}
	}

	// 4. RPC OrderDetail → types.OrderDetail（字段一一对应，直接赋值）
	o := rpcResp.Order
	return &types.OrderDetail{
		Id:          o.Id,
		OrderNo:     o.OrderNo,
		ElderlyId:   o.ElderlyId,
		NurseId:     o.NurseId,
		ServiceType: o.ServiceType,
		Status:      o.Status,
		Urgency:     o.Urgency,
		Amount:      o.Amount,
		ScheduledAt: o.ScheduledAt,
		StartedAt:   o.StartedAt,
		FinishedAt:  o.FinishedAt,
		CreatedAt:   o.CreatedAt,
	}, nil
}
