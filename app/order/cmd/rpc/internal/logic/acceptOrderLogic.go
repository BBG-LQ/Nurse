package logic

import (
	"context"
	"errors"

	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAcceptOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptOrderLogic {
	return &AcceptOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AcceptOrderLogic) AcceptOrder(in *order.AcceptOrderReq) (*order.AcceptOrderResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验

	// 2. 查询订单
	orderInfo, err := l.svcCtx.ServiceOrderModel.FindOne(l.ctx, uint64(in.OrderId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "订单不存在")
		}
		l.Errorf("查询订单失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 3. 幂等：已确认接单 且 归属该护士 → 直接成功
	if orderInfo.Status == orderStatusAccepted &&
		orderInfo.NurseId.Valid && orderInfo.NurseId.Int64 == in.NurseId {
		return &order.AcceptOrderResp{Ok: true}, nil
	}

	// 4. 订单尚未派单：没有护士可确认
	if !orderInfo.NurseId.Valid || orderInfo.NurseId.Int64 <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单尚未派单，无法确认")
	}

	// 5. 归属校验：只能确认派给自己的订单
	if orderInfo.NurseId.Int64 != in.NurseId {
		l.Infof("护士尝试确认他人订单: orderId=%d, orderNurseId=%d, reqNurseId=%d",
			in.OrderId, orderInfo.NurseId.Int64, in.NurseId)
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权确认该订单")
	}

	// 6. 状态校验：只有「已派待确认」态可确认
	if orderInfo.Status != orderStatusPending {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "当前订单状态不允许确认")
	}

	// 7. 乐观锁推进状态：0 → 1
	orderInfo.Status = orderStatusAccepted

	err = l.svcCtx.ServiceOrderModel.UpdateWithVersion(l.ctx, nil, orderInfo)
	if err != nil {
		if errors.Is(err, model.ErrNoRowsUpdate) {
			// 版本冲突：并发确认 / 用户取消等改动了订单
			l.Infof("确认接单版本冲突: orderId=%d, nurseId=%d", in.OrderId, in.NurseId)
			return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "订单状态已变更，请刷新后重试")
		}
		l.Errorf("确认接单更新失败: orderId=%d, nurseId=%d, err=%v",
			in.OrderId, in.NurseId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	l.Infof("护士确认接单成功: orderId=%d, nurseId=%d", in.OrderId, in.NurseId)
	return &order.AcceptOrderResp{Ok: true}, nil

}
