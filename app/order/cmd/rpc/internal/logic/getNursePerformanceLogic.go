package logic

import (
	"context"

	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetNursePerformanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	maxInt32 = int64(^uint32(0) >> 1) // 2147483647
)

func NewGetNursePerformanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNursePerformanceLogic {
	return &GetNursePerformanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNursePerformanceLogic) GetNursePerformance(in *order.GetNursePerformanceReq) (*order.GetNursePerformanceResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验

	// 2. 越权校验：护士只能查看自己的绩效
	//    若该接口还供管理端使用，需在 ctx 注入角色后再放开此校验（见第五节）

	// 3. 统计订单数
	orderCount, err := l.svcCtx.ServiceOrderModel.FindCount(
		l.ctx,
		l.newFinishedOrderBuilder(in.NurseId),
		"id",
	)
	if err != nil {
		l.Errorf("统计护士订单数失败: nurseId=%d, err=%v", in.NurseId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 4. 统计服务总时长
	//    FindSum 内部拼为 IFNULL(SUM(<field>),0)，传表达式即可得到 SUM(TIMESTAMPDIFF(...))
	//    注意：squirrel builder 是值类型，Columns() 返回新副本，但为清晰起见每次重建
	minutes, err := l.svcCtx.ServiceOrderModel.FindSum(
		l.ctx,
		l.newFinishedOrderBuilder(in.NurseId),
		"TIMESTAMPDIFF(MINUTE, started_at, finished_at)",
	)
	if err != nil {
		l.Errorf("统计护士服务时长失败: nurseId=%d, err=%v", in.NurseId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 5. 计算评分均值
	scoreAvg, err := l.calcScoreAvg(in.NurseId)
	if err != nil {
		return nil, err
	}

	return &order.GetNursePerformanceResp{
		OrderCount:     clampInt64ToInt32(orderCount),
		ServiceMinutes: clampFloatToInt32(minutes),
		ScoreAvg:       scoreAvg,
	}, nil
}

// newFinishedOrderBuilder 每次构造新 builder，避免被 FindCount/FindSum 内部 Columns() 影响
func (l *GetNursePerformanceLogic) newFinishedOrderBuilder(nurseId int64) squirrel.SelectBuilder {
	return l.svcCtx.ServiceOrderModel.SelectBuilder().Where(squirrel.Eq{
		"nurse_id": nurseId,
		"status":   orderStatusFinished,
	})
}

// calcScoreAvg 从 satisfaction 表实时计算护士评分均值；无评价时返回 0
func (l *GetNursePerformanceLogic) calcScoreAvg(nurseId int64) (float64, error) {
	count, err := l.svcCtx.SatisfactionModel.FindCount(
		l.ctx,
		l.svcCtx.SatisfactionModel.SelectBuilder().Where(squirrel.Eq{"nurse_id": nurseId}),
		"id",
	)
	if err != nil {
		l.Errorf("统计护士评价数失败: nurseId=%d, err=%v", nurseId, err)
		return 0, xerr.NewErrCode(xerr.DbError)
	}
	if count == 0 {
		return 0, nil
	}

	sum, err := l.svcCtx.SatisfactionModel.FindSum(
		l.ctx,
		l.svcCtx.SatisfactionModel.SelectBuilder().Where(squirrel.Eq{"nurse_id": nurseId}),
		"score",
	)
	if err != nil {
		l.Errorf("统计护士评分总分失败: nurseId=%d, err=%v", nurseId, err)
		return 0, xerr.NewErrCode(xerr.DbError)
	}

	return sum / float64(count), nil
}

// clampInt64ToInt32 防止 int64 → int32 溢出（理论不可能，但避免静默截断）
func clampInt64ToInt32(v int64) int32 {
	if v <= 0 {
		return 0
	}
	if v > maxInt32 {
		return int32(maxInt32)
	}
	return int32(v)
}

// clampFloatToInt32 处理负数（异常数据如 finished_at < started_at）与溢出
func clampFloatToInt32(v float64) int32 {
	if v <= 0 {
		return 0
	}
	if v > float64(maxInt32) {
		return int32(maxInt32)
	}
	return int32(v)
}
