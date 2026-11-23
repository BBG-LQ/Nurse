package logic

import (
	"context"
	"database/sql"
	"errors"

	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// 时间输出格式，与 CreateOrder 的 scheduledAt 解析格式保持一致
const orderTimeLayout = "2006-01-02 15:04:05"

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.GetOrderResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验

	// 2. 查询订单（FindOne 走缓存，miss 时回源 DB）
	m, err := l.svcCtx.ServiceOrderModel.FindOne(l.ctx, uint64(in.OrderId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "订单不存在")
		}
		l.Errorf("查询订单失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 3. 组装 OrderDetail
	detail := &order.OrderDetail{
		Id:          int64(m.Id),
		OrderNo:     m.OrderNo,
		ElderlyId:   int64(m.ElderlyId),
		NurseId:     m.NurseId.Int64, // 未派单时为 0（sql.NullInt64 零值 Int64 即 0）
		ServiceType: m.ServiceType,
		Status:      int32(m.Status),
		Urgency:     int32(m.Urgency),
		Amount:      m.Amount.Float64,
		ScheduledAt: formatNullTime(m.ScheduledAt),
		StartedAt:   formatNullTime(m.StartedAt),
		FinishedAt:  formatNullTime(m.FinishedAt),
		CreatedAt:   m.CreateTime.Format(orderTimeLayout),
	}

	return &order.GetOrderResp{Order: detail}, nil
}

// formatNullTime 将 sql.NullTime 格式化为字符串，无效值返回空串
func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(orderTimeLayout)
}
