package logic

import (
	"context"
	"database/sql"
	"time"

	"Nurse/app/vital/cmd/rpc/internal/svc"
	"Nurse/app/vital/cmd/rpc/vital"
	"Nurse/app/vital/model"
	"Nurse/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
)

type IngestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIngestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IngestLogic {
	return &IngestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IngestLogic) Ingest(in *vital.IngestReq) (*vital.IngestResp, error) {
	// todo: add your logic here and delete this line

	abnormal := l.checkAbnormal(in)

	record := &model.VitalRecord{
		ElderlyId:  uint64(in.ElderlyId),
		HeartRate:  sql.NullInt64{Int64: int64(in.HeartRate), Valid: in.HeartRate > 0},
		Systolic:   sql.NullInt64{Int64: int64(in.Systolic), Valid: in.Systolic > 0},
		Diastolic:  sql.NullInt64{Int64: int64(in.Diastolic), Valid: in.Diastolic > 0},
		Spo2:       sql.NullFloat64{Float64: in.Spo2, Valid: in.Spo2 > 0},
		Posture:    sql.NullInt64{Int64: int64(in.Posture), Valid: true},
		Source:     "device",
		RecordTime: time.Now(),
		Status:     1, // 1可接单
		DelState:   globalkey.DelStateNo,
		Version:    0,
	}

	result, err := l.svcCtx.VitalRecordModel.Insert(l.ctx, nil, record)
	if err != nil {
		l.Errorf("Ingest insert vital record failed: %v", err)
		return nil, err
	}

	recordId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("Ingest get last insert id failed: %v", err)
		return nil, err
	}

	return &vital.IngestResp{

		Abnormal: abnormal,
		RecordId: recordId,
	}, nil
}

// checkAbnormal 根据生命体征指标判定是否异常
func (l *IngestLogic) checkAbnormal(in *vital.IngestReq) bool {
	// 心率：正常 60~100，<50 或 >120 视为异常
	if in.HeartRate > 0 && (in.HeartRate < 50 || in.HeartRate > 120) {
		return true
	}
	// 收缩压：正常 90~120，<90 或 >140 视为异常
	if in.Systolic > 0 && (in.Systolic < 90 || in.Systolic > 140) {
		return true
	}
	// 舒张压：正常 60~80，<60 或 >90 视为异常
	if in.Diastolic > 0 && (in.Diastolic < 60 || in.Diastolic > 90) {
		return true
	}
	// 血氧饱和度：正常 >=95，低于则异常
	if in.Spo2 > 0 && in.Spo2 < 95 {
		return true
	}
	// 跌倒（posture=3）是明确异常信号
	if in.Posture == 3 {
		return true
	}
	return false
}
