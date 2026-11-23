package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	// 分布式锁
	dispatchLockPrefix = "lock:order:dispatch:"
	dispatchLockExpire = 10 // 秒

	// 每个订单最多尝试 N 个护士（用于处理负载被并发抢光的场景）
	dispatchMaxRetries = 3

	// 派单时护士最大可承受负载（超过后不再分配）
	dispatchDefaultMaxNurseLoad = 10

	// 订单状态
	orderStatusPending  = 0 // 待派
	orderStatusAccepted = 1 // 已接单
)

type DispatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchLogic {
	return &DispatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DispatchLogic) Dispatch(in *order.DispatchReq) (*order.DispatchResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验

	// 2. 分布式锁：防止同一订单被并发重复派单
	lockKey := fmt.Sprintf("%s%d", dispatchLockPrefix, in.OrderId)
	locked, err := l.svcCtx.RedisClient.SetnxEx(lockKey, "1", dispatchLockExpire)
	if err != nil {
		l.Errorf("获取派单锁失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}
	if !locked {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单正在派单中，请稍后重试")
	}
	defer func() {
		if _, err := l.svcCtx.RedisClient.Del(lockKey); err != nil {
			l.Errorf("释放派单锁失败: orderId=%d, err=%v", in.OrderId, err)
		}
	}()

	// 3. 查询订单
	orderInfo, err := l.svcCtx.ServiceOrderModel.FindOne(l.ctx, uint64(in.OrderId))
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "订单不存在")
		}
		l.Errorf("查询订单失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 4. 幂等：已派过单，直接返回
	if orderInfo.NurseId.Valid && orderInfo.NurseId.Int64 > 0 {
		return &order.DispatchResp{NurseId: orderInfo.NurseId.Int64}, nil
	}

	// 5. 状态校验：只有待派订单可以派单
	if orderInfo.Status != orderStatusPending {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "当前订单状态不允许派单")
	}

	// 6. 查询老人信息（校验老人存在；如需按老人区域匹配护士，可用返回值扩展）
	if _, err := l.svcCtx.ElderlyRpc.GetElderly(l.ctx, &elderly.GetElderlyReq{
		Id: int64(orderInfo.ElderlyId),
	}); err != nil {
		l.Errorf("查询老人信息失败: elderlyId=%d, err=%v", orderInfo.ElderlyId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 7. 选择护士并派单，最多尝试 dispatchMaxRetries 个护士
	triedNurseIds := make(map[int64]struct{})
	for attempt := 0; attempt < dispatchMaxRetries; attempt++ {
		// 7.1 查询可派单护士（FindAll 已按 current_load ASC, score_avg DESC 排序）
		nurseList, err := l.svcCtx.NurseRpc.ListDispatchableNurses(l.ctx, &nurse.ListDispatchableNursesReq{
			ServiceType: orderInfo.ServiceType,
			MaxLoad:     dispatchDefaultMaxNurseLoad,
		})
		if err != nil {
			l.Errorf("查询可派单护士失败: orderId=%d, serviceType=%s, err=%v",
				in.OrderId, orderInfo.ServiceType, err)
			return nil, xerr.NewErrCode(xerr.DbError)
		}

		// 7.2 选出本轮未被尝试过的第一个护士
		var selected *nurse.NurseBrief
		for _, n := range nurseList.Nurses {
			if _, tried := triedNurseIds[n.Id]; tried {
				continue
			}
			selected = n
			break
		}
		if selected == nil {
			return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "暂无可用护士")
		}

		// 7.3 预占护士负载（+1）
		if _, err = l.svcCtx.NurseRpc.IncrLoad(l.ctx, &nurse.IncrLoadReq{
			NurseId: selected.Id,
			Delta:   1,
		}); err != nil {
			// 该护士负载可能刚好被并发占满，换下一个
			triedNurseIds[selected.Id] = struct{}{}
			l.Infof("护士负载预占失败，尝试下一个: orderId=%d, nurseId=%d, err=%v",
				in.OrderId, selected.Id, err)
			continue
		}

		// 7.4 更新订单前重新拉取（保证 version 是最新的，避免乐观锁失败）
		orderInfo, err = l.svcCtx.ServiceOrderModel.FindOne(l.ctx, uint64(in.OrderId))
		if err != nil {
			l.Errorf("重新查询订单失败: orderId=%d, err=%v", in.OrderId, err)
			l.releaseNurseLoad(selected.Id)
			return nil, xerr.NewErrCode(xerr.DbError)
		}
		// 双检：期间订单状态若被变更，则回滚护士负载并返回
		if orderInfo.NurseId.Valid && orderInfo.NurseId.Int64 > 0 {
			l.releaseNurseLoad(selected.Id)
			return &order.DispatchResp{NurseId: orderInfo.NurseId.Int64}, nil
		}
		if orderInfo.Status != orderStatusPending {
			l.releaseNurseLoad(selected.Id)
			return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单状态已变更，无法派单")
		}

		// 7.5 事务更新订单（乐观锁）
		orderInfo.NurseId = sql.NullInt64{Int64: selected.Id, Valid: true}
		orderInfo.Status = orderStatusAccepted

		updateErr := l.svcCtx.ServiceOrderModel.Trans(l.ctx,
			func(ctx context.Context, session sqlx.Session) error {
				return l.svcCtx.ServiceOrderModel.UpdateWithVersion(ctx, session, orderInfo)
			},
		)
		if updateErr == nil {
			l.Infof("派单成功: orderId=%d, nurseId=%d", in.OrderId, selected.Id)
			return &order.DispatchResp{NurseId: selected.Id}, nil
		}

		// 7.6 订单更新失败，补偿回滚护士负载
		l.Errorf("更新订单失败，回滚护士负载: orderId=%d, nurseId=%d, err=%v",
			in.OrderId, selected.Id, updateErr)
		l.releaseNurseLoad(selected.Id)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 8. 重试耗尽
	return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "派单失败，请重试")
}

// releaseNurseLoad 回滚护士负载（-1），失败仅记录日志（不阻塞主流程）
func (l *DispatchLogic) releaseNurseLoad(nurseId int64) {
	if _, err := l.svcCtx.NurseRpc.IncrLoad(l.ctx, &nurse.IncrLoadReq{
		NurseId: nurseId,
		Delta:   -1,
	}); err != nil {
		l.Errorf("回滚护士负载失败: nurseId=%d, err=%v", nurseId, err)
	}
}
