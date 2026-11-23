package logic

import (
	"context"
	"database/sql"
	"time"

	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/app/alter/cmd/rpc/internal/svc"
	"Nurse/app/alter/model"
	"Nurse/pkg/globalkey"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SosNurseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSosNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SosNurseLogic {
	return &SosNurseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SosNurseLogic) SosNurse(in *alert.SosNurseReq) (*alert.SosNurseResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不能为空")
	}
	if in.OrderId < 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}

	// if in.OrderId > 0 {
	// 	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, uint64(in.OrderId))
	// 	if err == nil && order.ElderlyId.Valid {
	// 		alertEvent.ElderlyId = order.ElderlyId
	// 	}
	// }

	// 2. 构造高危 SOS 报警记录
	//    - AlertType 固定 sos_nurse，Severity 固定 2（高危）
	//    - 请求未带 elderly_id，若 order_id > 0 可反查补齐（见文末增强项）
	alertEvent := &model.AlertEvent{
		NurseId:   sql.NullInt64{Int64: in.NurseId, Valid: true},
		OrderId:   sql.NullInt64{Int64: in.OrderId, Valid: in.OrderId > 0},
		AlertType: "sos_nurse",
		Severity:  2, // SOS 固定高危
		Content:   sql.NullString{String: in.Content, Valid: in.Content != ""},
		Handled:   0,
		DelState:  globalkey.DelStateNo,
		Version:   1,
	}

	// 3. 落库
	result, err := l.svcCtx.AlertEventModel.Insert(l.ctx, nil, alertEvent)
	if err != nil {
		l.Errorf("sos nurse insert failed, nurse_id=%d, order_id=%d, err=%v", in.NurseId, in.OrderId, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "insert sos nurse alert failed: %v", err)
	}

	alertId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("sos nurse get last insert id failed, err=%v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "get last insert id failed: %v", err)
	}

	// 4. TODO: 高危报警触发实时通知（WebSocket 推管理台 / 调度中心短信 / 派单）
	//    建议异步执行，避免阻塞护士端 SOS 响应
	l.Infof("nurse SOS alert created, alert_id=%d, nurse_id=%d, order_id=%d", alertId, in.NurseId, in.OrderId)

	builder := l.svcCtx.AlertEventModel.SelectBuilder().
		Where("nurse_id = ?", in.NurseId).
		Where("alert_type = ?", "sos_nurse").
		Where("handled = ?", 0).
		Where("create_time >= ?", time.Now().Add(-30*time.Second))
	recent, _ := l.svcCtx.AlertEventModel.FindAll(l.ctx, builder, "id DESC")
	if len(recent) > 0 {
		return &alert.SosNurseResp{AlertId: int64(recent[0].Id)}, nil
	}

	return &alert.SosNurseResp{AlertId: alertId}, nil
}
