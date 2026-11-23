package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	orderStatusServing  = 2 // 服务中
	orderStatusFinished = 3 // 结束服务
)

type StartOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartOrderLogic {
	return &StartOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StartOrderLogic) StartOrder(in *order.StartOrderReq) (*order.StartOrderResp, error) {
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

	// 3. 幂等：已服务中 / 已完成，视为开始成功（重复点击或网络重试）
	if orderInfo.Status == orderStatusServing || orderInfo.Status == orderStatusFinished {
		return &order.StartOrderResp{Ok: true}, nil
	}

	// 4. 越权校验：只有该订单的接单护士能开始服务
	//    身份从 ctx 拿（鉴权中间件应已注入 jwtUserId）
	if uid := ctxdata.GetUidFromCtx(l.ctx); uid > 0 {
		if !orderInfo.NurseId.Valid || orderInfo.NurseId.Int64 != uid {
			l.Infof("护士尝试开始他人订单: orderId=%d, orderNurseId=%v, reqUid=%d",
				in.OrderId, orderInfo.NurseId, uid)
			return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权操作该订单")
		}
	}

	// 5. 状态校验：只有「已接单」可开始服务
	if orderInfo.Status != orderStatusAccepted {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "当前订单状态不允许开始服务")
	}

	// 6. 乐观锁推进 1 → 2，同时写 started_at
	orderInfo.Status = orderStatusServing
	orderInfo.StartedAt = sql.NullTime{Time: time.Now(), Valid: true}

	err = l.svcCtx.ServiceOrderModel.UpdateWithVersion(l.ctx, nil, orderInfo)
	if err != nil {
		if errors.Is(err, model.ErrNoRowsUpdate) {
			// 版本冲突：并发开始 / 期间被取消等
			l.Infof("开始服务版本冲突: orderId=%d, nurseId=%v", in.OrderId, orderInfo.NurseId)
			return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "订单状态已变更，请刷新后重试")
		}
		l.Errorf("开始服务更新失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	l.Infof("开始服务成功: orderId=%d, nurseId=%v, startedAt=%v",
		in.OrderId, orderInfo.NurseId, orderInfo.StartedAt.Time)
	return &order.StartOrderResp{Ok: true}, nil
}
