package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/globalkey"
	"Nurse/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	satisfactionScoreMin   = 1
	satisfactionScoreMax   = 5
	satisfactionCommentMax = 500 // 评论最大字符数
)

type SubmitSatisfactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitSatisfactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitSatisfactionLogic {
	return &SubmitSatisfactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 提交满意度评价
func (l *SubmitSatisfactionLogic) SubmitSatisfaction(in *order.SubmitSatisfactionReq) (*order.SubmitSatisfactionResp, error) {
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

	// 3. 订单状态校验：只有已完成订单可评价
	if orderInfo.Status != orderStatusFinished {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单尚未完成，无法评价")
	}
	if !orderInfo.NurseId.Valid || orderInfo.NurseId.Int64 <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单未分配护士")
	}

	// 4. 护士一致性校验：req.NurseId 若传了必须与订单一致，防止评错人
	// if in.NurseId > 0 && in.NurseId != orderInfo.NurseId.Int64 {
	// 	l.Infof("护士ID不一致: orderId=%d, orderNurseId=%d, reqNurseId=%d",
	// 		in.OrderId, orderInfo.NurseId.Int64, in.NurseId)
	// 	return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID与订单不符")
	// }

	// 5. 身份校验：评价人必须是订单所属老人
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.TokenExpireError, "未登录或登录已过期")
	}
	if uint64(uid) != orderInfo.ElderlyId {
		l.Infof("非本人评价: orderId=%d, orderElderlyId=%d, reqUid=%d",
			in.OrderId, orderInfo.ElderlyId, uid)
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权评价该订单")
	}

	// 6. 幂等：已评价则直接返回成功（一单一评，重复提交不报错）
	if _, err := l.svcCtx.SatisfactionModel.FindOneByOrderId(l.ctx, uint64(in.OrderId)); err == nil {
		l.Infof("订单已评价，幂等返回: orderId=%d", in.OrderId)
		return &order.SubmitSatisfactionResp{Ok: true}, nil
	} else if !errors.Is(err, model.ErrNotFound) {
		l.Errorf("查询评价记录失败: orderId=%d, err=%v", in.OrderId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 7. 组装评价数据
	//    DeleteTime / DelState / Version 在 model_gen 的 Insert 里被注释掉了，
	//    必须由调用方显式设置（与 CreateOrder 相同）
	comment := strings.TrimSpace(in.Comment)
	data := &model.Satisfaction{
		OrderId:    uint64(in.OrderId),
		NurseId:    uint64(orderInfo.NurseId.Int64), // 以订单为准，不信任客户端
		Score:      int64(in.Score),
		Comment:    sql.NullString{String: comment, Valid: comment != ""},
		DeleteTime: time.Unix(0, 0),
		DelState:   globalkey.DelStateNo,
		Version:    0,
	}

	// 8. 插入
	if _, err := l.svcCtx.SatisfactionModel.Insert(l.ctx, nil, data); err != nil {
		// 并发重复插入兜底：重新查询若命中，说明是并发重复提交，视为幂等成功
		if _, findErr := l.svcCtx.SatisfactionModel.FindOneByOrderId(l.ctx, uint64(in.OrderId)); findErr == nil {
			l.Infof("并发重复提交，幂等返回: orderId=%d", in.OrderId)
			return &order.SubmitSatisfactionResp{Ok: true}, nil
		}
		l.Errorf("插入评价失败: orderId=%d, nurseId=%d, err=%v",
			in.OrderId, data.NurseId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 9. 刷新护士评分：从 satisfaction 聚合最新均分，调 nurse-rpc.UpdateScore
	//    评价已落库，聚合失败只告警（由定时任务兜底），不能回滚评价
	//    刷新的是 NurseProfile.ScoreAvg，供 ListDispatchableNurses 派单排序使用
	if err := l.recalcNurseScoreAvg(int64(data.NurseId)); err != nil {
		l.Errorf("刷新护士评分失败[需对账]: orderId=%d, nurseId=%d, err=%v",
			in.OrderId, data.NurseId, err)
	}

	l.Infof("评价提交成功: orderId=%d, nurseId=%d, score=%d, hasComment=%v",
		in.OrderId, data.NurseId, in.Score, data.Comment.Valid)
	return &order.SubmitSatisfactionResp{Ok: true}, nil
}

// recalcNurseScoreAvg 从 satisfaction 表聚合护士最新评分均值，并回写 nurse_profile
// 返回错误仅用于调用方记录日志，不影响评价提交结果
func (l *SubmitSatisfactionLogic) recalcNurseScoreAvg(nurseId int64) error {
	// 9.1 统计评价条数
	count, err := l.svcCtx.SatisfactionModel.FindCount(
		l.ctx,
		l.svcCtx.SatisfactionModel.SelectBuilder().Where(squirrel.Eq{"nurse_id": nurseId}),
		"id",
	)
	if err != nil {
		return fmt.Errorf("统计护士评价数失败: %w", err)
	}
	// 理论不会为 0（刚插入了一条），防御性处理避免除零
	if count == 0 {
		return nil
	}

	// 9.2 统计评分总分
	sum, err := l.svcCtx.SatisfactionModel.FindSum(
		l.ctx,
		l.svcCtx.SatisfactionModel.SelectBuilder().Where(squirrel.Eq{"nurse_id": nurseId}),
		"score",
	)
	if err != nil {
		return fmt.Errorf("统计护士评分总分失败: %w", err)
	}

	avg := sum / float64(count)

	// 9.3 回写护士档案评分
	if _, err := l.svcCtx.NurseRpc.UpdateScore(l.ctx, &nurse.UpdateScoreReq{
		NurseId:  nurseId,
		ScoreAvg: avg,
	}); err != nil {
		return fmt.Errorf("调用 UpdateScore 失败: %w", err)
	}

	l.Infof("护士评分已刷新: nurseId=%d, count=%d, avg=%.4f", nurseId, count, avg)
	return nil
}
