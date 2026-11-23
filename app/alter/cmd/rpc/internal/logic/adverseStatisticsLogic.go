package logic

import (
	"context"
	"time"

	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/app/alter/cmd/rpc/internal/svc"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

// 合法的不良事件类型，用于统计时预置 0 值，保证前端图表分类轴稳定
var adverseEventTypes = []string{"fall", "pressure_ulcer", "medication_error", "dispute"}

type AdverseStatisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdverseStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdverseStatisticsLogic {
	return &AdverseStatisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdverseStatisticsLogic) AdverseStatistics(in *alert.AdverseStatisticsReq) (*alert.AdverseStatisticsResp, error) {
	// todo: add your logic here and delete this line

	// 1. 解析月份，空则默认当月
	var start time.Time
	if in.Month == "" {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	} else {
		t, err := time.ParseInLocation("2006-01", in.Month, time.Local)
		if err != nil {
			return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "月份格式错误，应为 YYYY-MM")
		}
		start = t
	}
	end := start.AddDate(0, 1, 0) // 次月 1 号，区间左闭右开

	// 2. 查询当月全部未删除的不良事件（FindAll 内部已叠加 del_state 过滤）
	builder := l.svcCtx.AdverseEventModel.SelectBuilder().
		Where("create_time >= ?", start).
		Where("create_time < ?", end)

	rows, err := l.svcCtx.AdverseEventModel.FindAll(l.ctx, builder, "id DESC")
	if err != nil {
		l.Errorf("adverse statistics query failed, month=%s, err=%v", in.Month, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "query adverse events failed: %v", err)
	}

	// 3. 内存聚合：先预置合法类型为 0，保证前端分类轴完整
	typeCount := make(map[string]int32, len(adverseEventTypes))
	for _, t := range adverseEventTypes {
		typeCount[t] = 0
	}
	for _, row := range rows {
		typeCount[row.EventType]++
	}

	return &alert.AdverseStatisticsResp{
		TypeCount: typeCount,
		Total:     int32(len(rows)),
	}, nil
}
