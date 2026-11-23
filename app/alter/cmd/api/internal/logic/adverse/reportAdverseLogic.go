// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adverse

import (
	"context"
	"fmt"

	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/tool"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const reporterSalt = "nurse-adverse-report-salt-v1" // 建议迁到 config

type ReportAdverseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 匿名上报
func NewReportAdverseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportAdverseLogic {
	return &ReportAdverseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportAdverseLogic) ReportAdverse(req *types.ReportAdverseReq) (resp *types.ReportAdverseResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 校验
	switch req.EventType {
	case "fall", "pressure_ulcer", "medication_error", "dispute":
	default:
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "事件类型不合法")
	}
	if req.Description == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "事件描述不能为空")
	}

	// 2. 确定 nurseId：显式指定 > JWT 当前用户 > 匿名(0)
	nurseId := req.NurseId
	if nurseId <= 0 {
		nurseId = ctxdata.GetUidFromCtx(l.ctx)
	}

	// 3. 匿名时生成 reporter_hash（实名则留空）
	var reporterHash string
	if nurseId <= 0 {
		ip := ctxdata.GetRemoteIPFromCtx(l.ctx)
		if ip != "" {
			reporterHash = tool.Md5ByString(fmt.Sprintf("%s|%s", ip, reporterSalt))
		}
	}

	// 4. 转发 RPC
	rpcResp, err := l.svcCtx.AlertRpc.ReportAdverse(l.ctx, &alert.ReportAdverseReq{
		EventType:    req.EventType,
		Description:  req.Description,
		NurseId:      nurseId,
		ElderlyId:    req.ElderlyId,
		ReporterHash: reporterHash,
	})
	if err != nil {
		l.Errorf("call rpc ReportAdverse failed, req=%+v, err=%v", req, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "report adverse failed: %v", err)
	}

	return &types.ReportAdverseResp{Id: rpcResp.Id}, nil
}
