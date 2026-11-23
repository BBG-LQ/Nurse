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

type ReportAdverseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportAdverseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportAdverseLogic {
	return &ReportAdverseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportAdverseLogic) ReportAdverse(in *alert.ReportAdverseReq) (*alert.ReportAdverseResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	switch in.EventType {
	case "fall", "pressure_ulcer", "medication_error", "dispute":
		// 合法类型
	default:
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "事件类型不合法")
	}
	if in.Description == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "事件描述不能为空")
	}
	if in.NurseId <= 0 && in.ReporterHash == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不能为空")
	}
	if in.ElderlyId < 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "老人ID不合法")
	}

	// 2. 匿名判定：带 reporter_hash 即匿名上报
	anonymous := int64(0)
	if in.ReporterHash != "" {
		anonymous = 1
	}

	// 3. 构造不良事件记录
	//    匿名时清空 nurse_id，避免通过关联关系反推举报人
	adverseEvent := &model.AdverseEvent{
		EventType:    in.EventType,
		Description:  sql.NullString{String: in.Description, Valid: true},
		Anonymous:    anonymous,
		ReporterHash: sql.NullString{String: in.ReporterHash, Valid: in.ReporterHash != ""},
		ElderlyId:    sql.NullInt64{Int64: in.ElderlyId, Valid: in.ElderlyId > 0},
		Status:       0, // 0待复盘
		DelState:     globalkey.DelStateNo,
		Version:      1,
	}
	if anonymous == 0 {
		adverseEvent.NurseId = sql.NullInt64{Int64: in.NurseId, Valid: true}
	} else {
		adverseEvent.NurseId = sql.NullInt64{Valid: false} // 匿名字段置 NULL
	}

	// 4. 落库
	result, err := l.svcCtx.AdverseEventModel.Insert(l.ctx, nil, adverseEvent)
	if err != nil {
		l.Errorf("report adverse insert failed, req=%+v, err=%v", in, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "insert adverse event failed: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("report adverse get last insert id failed, err=%v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "get last insert id failed: %v", err)
	}

	l.Infof("adverse event reported, id=%d, event_type=%s, anonymous=%d", id, in.EventType, anonymous)

	return &alert.ReportAdverseResp{Id: id}, nil
}
