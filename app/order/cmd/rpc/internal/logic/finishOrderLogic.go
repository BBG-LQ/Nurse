package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type FinishOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFinishOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinishOrderLogic {
	return &FinishOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FinishOrderLogic) FinishOrder(in *order.FinishOrderReq) (*order.FinishOrderResp, error) {
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

	// 3. 幂等：已完成 → 直接成功（重复点击/网络重试）
	//    注意：不重复释放负载，符合幂等语义
	if orderInfo.Status == orderStatusFinished {
		return &order.FinishOrderResp{Ok: true}, nil
	}

	// 4. 越权校验：只有该订单的接单护士能完成服务
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.TokenExpireError, "未登录或登录已过期")
	}
	if !orderInfo.NurseId.Valid || orderInfo.NurseId.Int64 != uid {
		l.Infof("护士尝试完成他人订单: orderId=%d, orderNurseId=%v, reqUid=%d",
			in.OrderId, orderInfo.NurseId, uid)
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权操作该订单")
	}

	// 5. 状态校验：只有「服务中」可完成
	if orderInfo.Status != orderStatusServing {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "当前订单状态不允许完成")
	}

	// 6. 乐观锁推进 2 → 3，写 finished_at
	orderInfo.Status = orderStatusFinished
	orderInfo.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}

	err = l.svcCtx.ServiceOrderModel.UpdateWithVersion(l.ctx, nil, orderInfo)
	if err != nil {
		if errors.Is(err, model.ErrNoRowsUpdate) {
			l.Infof("完成服务版本冲突: orderId=%d, nurseId=%d", in.OrderId, orderInfo.NurseId.Int64)
			return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "订单状态已变更，请刷新后重试")
		}
		l.Errorf("完成服务更新订单失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 7. 释放护士负载（订单已落库，此处失败只能告警，不能回滚订单）
	//    统一走 DecrLoad：语义更清晰，与 Dispatch 的 IncrLoad(+1) 形成对称
	nurseId := orderInfo.NurseId.Int64
	if _, err := l.svcCtx.NurseRpc.DecrLoad(l.ctx, &nurse.DecrLoadReq{
		NurseId: nurseId,
		Delta:   1,
	}); err != nil {
		// 订单已完成，负载释放失败 → 护士负载虚高，影响后续派单
		// 记录 ERROR + metric，依赖对账任务修复；不能回滚订单状态
		l.Errorf("完成服务后释放护士负载失败[需对账]: orderId=%d, nurseId=%d, err=%v",
			in.OrderId, nurseId, err)
	}

	l.Infof("完成服务成功: orderId=%d, nurseId=%d, startedAt=%v, finishedAt=%v",
		in.OrderId, nurseId, orderInfo.StartedAt.Time, orderInfo.FinishedAt.Time)
	return &order.FinishOrderResp{Ok: true}, nil
}
