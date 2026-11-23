package logic

import (
	"context"
	"time"

	"Nurse/app/vital/cmd/rpc/internal/svc"
	"Nurse/app/vital/cmd/rpc/vital"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecordsLogic {
	return &ListRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListRecords 按老人 ID 查询最近 N 天的生命体征记录，默认 7 天，最多 90 天
func (l *ListRecordsLogic) ListRecords(in *vital.ListRecordsReq) (*vital.ListRecordsResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "elderly_id 不能为空")
	}

	days := int(in.Days)
	if days <= 0 {
		days = 7 // 默认近 7 天
	}
	if days > 90 {
		days = 90 // 上限 90 天，防止全表扫描
	}

	since := time.Now().AddDate(0, 0, -days)

	// 2. 构建查询：elderly_id = ? AND record_time >= ?，按时间倒序
	//    del_state = ? 由 FindAll 内部统一追加，无需重复条件
	builder := l.svcCtx.VitalRecordModel.SelectBuilder().
		Where("elderly_id = ?", in.ElderlyId).
		Where("record_time >= ?", since)

	records, err := l.svcCtx.VitalRecordModel.FindAll(l.ctx, builder, "record_time DESC")
	if err != nil {
		l.Errorf("ListRecords query failed, elderlyId=%d, days=%d, err=%v", in.ElderlyId, days, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 3. 转换为 proto 返回（空结果返回空切片而非 nil）
	list := make([]*vital.VitalRecord, 0, len(records))
	for _, r := range records {

		temp := &vital.VitalRecord{
			Id: int64(r.Id),

			HeartRate: int32(r.HeartRate.Int64),
			Systolic:  int32(r.Systolic.Int64),
			Diastolic: int32(r.Diastolic.Int64),
			Spo2:      r.Spo2.Float64,
			Posture:   int32(r.Posture.Int64),

			RecordTime: r.RecordTime.Format(time.RFC3339),
		}
		list = append(list, temp)
	}

	return &vital.ListRecordsResp{
		List: list,
	}, nil
}
