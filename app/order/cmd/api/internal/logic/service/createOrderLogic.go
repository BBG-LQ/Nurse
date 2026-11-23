// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package service

import (
	"context"
	"unicode/utf8"

	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	urgencyNormal = 0
	urgencyUrgent = 1

	serviceTypeMaxLen = 64
	scheduledAtMaxLen = 32 // "2006-01-02 15:04:05" / RFC3339 均 < 32 字节
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建订单
func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验（快速失败，避免无谓 RPC 往返）
	if req.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "老人ID不合法")
	}
	if req.ServiceType == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "服务类型不能为空")
	}
	if utf8.RuneCountInString(req.ServiceType) > serviceTypeMaxLen {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "服务类型过长")
	}
	if req.Urgency != urgencyNormal && req.Urgency != urgencyUrgent {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "紧急程度不合法")
	}
	if req.Amount < 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "金额不能为负数")
	}
	// if utf8.RuneCountInString(req.ScheduledAt) > scheduledAtMaxLen {
	// 	return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "预约时间格式不合法")
	// }
	// ScheduledAt 的实际解析与合法性由 RPC 层统一处理（支持空串、多种 layout），
	// API 层只做长度拦截，避免长垃圾串打到 RPC

	// 2. 越权校验（老人端）：ctx 有 uid 时，只能为自己下单
	//    若接口同时供管理端调用（无 jwt uid），则信任 req.ElderlyId
	if uid := ctxdata.GetUidFromCtx(l.ctx); uid > 0 && uid != req.ElderlyId {
		l.Infof("老人尝试为他人下单: reqElderlyId=%d, reqUid=%d", req.ElderlyId, uid)
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权为他人下单")
	}

	// 3. 调用 RPC
	rpcResp, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderReq{
		ElderlyId:   req.ElderlyId,
		ServiceType: req.ServiceType,
		Urgency:     req.Urgency,
		Amount:      req.Amount,
		ScheduledAt: req.ScheduledAt,
	})
	if err != nil {
		// 不重造错误：RPC 层已把业务错误封成 xerr.CodeError，
		// 由全局 SetErrorHandler 统一转 JSON 响应
		l.Errorf("调用 CreateOrder RPC 失败: req=%+v, err=%v", req, err)
		return nil, err
	}

	// 4. 映射响应
	return &types.CreateOrderResp{
		OrderId: rpcResp.OrderId,
		OrderNo: rpcResp.OrderNo,
	}, nil

}
