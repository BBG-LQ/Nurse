// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ingest

import (
	"context"

	"Nurse/app/vital/cmd/api/internal/svc"
	"Nurse/app/vital/cmd/api/internal/types"
	"Nurse/app/vital/cmd/rpc/vital"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type IngestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模拟手环上报
func NewIngestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IngestLogic {
	return &IngestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IngestLogic) Ingest(req *types.IngestReq) (resp *types.IngestResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if req.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "elderlyId 不能为空")
	}
	if req.HeartRate <= 0 && req.Systolic <= 0 && req.Diastolic <= 0 && req.Spo2 <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "至少需要一项生命体征数据")
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.VitalRpc.Ingest(l.ctx, &vital.IngestReq{
		ElderlyId: req.ElderlyId,
		HeartRate: req.HeartRate,
		Systolic:  req.Systolic,
		Diastolic: req.Diastolic,
		Spo2:      req.Spo2,
		Posture:   req.Posture,
	})
	if err != nil {
		l.Errorf("rpc Ingest failed, elderlyId=%d, err=%v", req.ElderlyId, err)
		return nil, xerr.NewErrCode(xerr.ServerCommonError)
	}

	return &types.IngestResp{
		Abnormal: rpcResp.Abnormal,
		RecordId: rpcResp.RecordId,
	}, nil
}
