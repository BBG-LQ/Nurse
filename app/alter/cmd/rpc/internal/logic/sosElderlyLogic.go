package logic

import (
	"context"
	"database/sql"

	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/app/alter/cmd/rpc/internal/svc"
	"Nurse/app/alter/model"
	"Nurse/pkg/globalkey"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SosElderlyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSosElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SosElderlyLogic {
	return &SosElderlyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SosElderlyLogic) SosElderly(in *alert.SosElderlyReq) (*alert.SosElderlyResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "老人ID不能为空")
	}

	// 2. 构造高危 SOS 报警记录（类型/级别由业务语义固定，不接受前端传入）
	alertEvent := &model.AlertEvent{
		ElderlyId: sql.NullInt64{Int64: in.ElderlyId, Valid: true},
		AlertType: "sos_elderly",
		Severity:  2, // SOS 固定高危
		Content:   sql.NullString{String: in.Content, Valid: in.Content != ""},
		Handled:   0,
		DelState:  globalkey.DelStateNo,
		Version:   1,
	}

	// 3. 落库
	result, err := l.svcCtx.AlertEventModel.Insert(l.ctx, nil, alertEvent)
	if err != nil {
		l.Errorf("sos elderly insert failed, elderly_id=%d, err=%v", in.ElderlyId, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "insert sos alert failed: %v", err)
	}

	alertId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("sos elderly get last insert id failed, err=%v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "get last insert id failed: %v", err)
	}

	// 4. TODO: 高危报警触发实时通知（WebSocket 推护士站 / 短信 / 工单系统）
	//    建议异步：go func() { notify(alertId) }()，避免阻塞 SOS 响应
	l.Infof("SOS alert created, alert_id=%d, elderly_id=%d", alertId, in.ElderlyId)

	return &alert.SosElderlyResp{AlertId: alertId}, nil
}
