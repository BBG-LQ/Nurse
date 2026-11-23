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

type CreateAlertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAlertLogic {
	return &CreateAlertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAlertLogic) CreateAlert(in *alert.CreateAlertReq) (*alert.CreateAlertResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "老人ID不能为空")
	}
	switch in.AlertType {
	case "vital_abnormal", "sos_elderly", "sos_nurse":
		// 合法类型
	default:
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "报警类型不合法")
	}
	if in.Severity != 1 && in.Severity != 2 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "报警级别不合法，1=提示 2=高危")
	}

	// 2. 构造 AlertEvent 记录
	alertEvent := &model.AlertEvent{
		ElderlyId: sql.NullInt64{Int64: in.ElderlyId, Valid: true},
		AlertType: in.AlertType,
		Severity:  int64(in.Severity),
		Content:   sql.NullString{String: in.Content, Valid: in.Content != ""},
		Handled:   0, // 未处理
		DelState:  globalkey.DelStateNo,
		Version:   1, // 乐观锁初始版本号
	}

	// 3. 插入数据库
	result, err := l.svcCtx.AlertEventModel.Insert(l.ctx, nil, alertEvent)
	if err != nil {
		l.Errorf("create alert event failed, req=%+v, err=%v", in, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "insert alert event failed: %v", err)
	}

	// 4. 取回自增 ID
	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("get last insert id failed, err=%v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "get last insert id failed: %v", err)
	}

	return &alert.CreateAlertResp{Id: id}, nil
}
