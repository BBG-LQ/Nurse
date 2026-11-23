package logic

import (
	"context"
	"time"

	"Nurse/app/vital/cmd/rpc/internal/svc"
	"Nurse/app/vital/cmd/rpc/vital"
	"Nurse/app/vital/model"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckContinuousAbnormalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

const (
	// 连续异常预警阈值（天）
	continuousAbnormalThreshold = 3
	// 回溯查询窗口（天），需 >= threshold，防止漏判
	continuousAbnormalWindow = 30
)

func NewCheckContinuousAbnormalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckContinuousAbnormalLogic {
	return &CheckContinuousAbnormalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckContinuousAbnormal 检查老人是否连续多天出现生命体征异常
func (l *CheckContinuousAbnormalLogic) CheckContinuousAbnormal(in *vital.CheckContinuousAbnormalReq) (*vital.CheckContinuousAbnormalResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "elderly_id 不能为空")
	}

	now := time.Now()
	since := now.AddDate(0, 0, -continuousAbnormalWindow)

	// 2. 查询窗口内该老人的所有记录（del_state 由 FindAll 统一处理）
	builder := l.svcCtx.VitalRecordModel.SelectBuilder().
		Where("elderly_id = ?", in.ElderlyId).
		Where("record_time >= ?", since)

	records, err := l.svcCtx.VitalRecordModel.FindAll(l.ctx, builder, "record_time DESC")
	if err != nil {
		l.Errorf("CheckContinuousAbnormal query failed, elderlyId=%d, err=%v", in.ElderlyId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 3. 按天聚合：当天是否有记录、当天是否出现异常
	dailyHasRecord := make(map[string]bool)
	dailyAbnormal := make(map[string]bool)
	for _, r := range records {
		day := r.RecordTime.Format("2006-01-02")
		dailyHasRecord[day] = true
		if !dailyAbnormal[day] && isVitalRecordAbnormal(r) {
			dailyAbnormal[day] = true
		}
	}

	// 4. 确定起始日：今天有记录则从今天算；今天无记录则从昨天算（避免因当天数据未上报漏判）
	todayKey := now.Format("2006-01-02")
	startDay := now
	if !dailyHasRecord[todayKey] {
		startDay = now.AddDate(0, 0, -1)
	}

	// 5. 从起始日往前逐天统计连续异常天数，遇"无异常"或"无记录"即中断
	days := 0
	for i := 0; i < continuousAbnormalWindow; i++ {
		day := startDay.AddDate(0, 0, -i).Format("2006-01-02")
		if dailyAbnormal[day] {
			days++
			continue
		}
		break
	}

	return &vital.CheckContinuousAbnormalResp{
		Continuous: days >= continuousAbnormalThreshold,
		Days:       int32(days),
	}, nil
}

// isVitalRecordAbnormal 判定单条生命体征记录是否异常
// 规则与 Ingest 保持一致：心率/血压/血氧超阈值，或姿势为跌倒(3)
func isVitalRecordAbnormal(r *model.VitalRecord) bool {
	if r.HeartRate.Valid && (r.HeartRate.Int64 < 50 || r.HeartRate.Int64 > 120) {
		return true
	}
	if r.Systolic.Valid && (r.Systolic.Int64 < 90 || r.Systolic.Int64 > 140) {
		return true
	}
	if r.Diastolic.Valid && (r.Diastolic.Int64 < 60 || r.Diastolic.Int64 > 90) {
		return true
	}
	if r.Spo2.Valid && r.Spo2.Float64 < 95 {
		return true
	}
	if r.Posture.Valid && r.Posture.Int64 == 3 { // 3 跌倒
		return true
	}
	return false
}
