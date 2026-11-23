// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package query

import (
	"context"

	"Nurse/app/vital/cmd/api/internal/svc"
	"Nurse/app/vital/cmd/api/internal/types"
	"Nurse/app/vital/cmd/rpc/vital"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 近 N 天趋势
func NewListRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecordsLogic {
	return &ListRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRecordsLogic) ListRecords(req *types.ListRecordsReq) (resp *types.ListRecordsResp, err error) {
	// todo: add your logic here and delete this line
	// 1. 参数校验（Days 已有 default=7，无需再兜底，但做上限保护）
	if req.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "elderlyId 不能为空")
	}
	if req.Days <= 0 {
		req.Days = 7
	}
	if req.Days > 90 {
		req.Days = 90 // 与 RPC 层上限保持一致
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.VitalRpc.ListRecords(l.ctx, &vital.ListRecordsReq{
		ElderlyId: req.ElderlyId,
		Days:      req.Days,
	})
	if err != nil {
		l.Errorf("rpc ListRecords failed, elderlyId=%d, days=%d, err=%v", req.ElderlyId, req.Days, err)
		return nil, xerr.NewErrCode(xerr.ServerCommonError)
	}

	// 3. 转换响应（空结果返回空切片，避免前端拿到 null）
	list := make([]types.VitalRecordItem, 0, len(rpcResp.List))
	for _, r := range rpcResp.List {
		list = append(list, types.VitalRecordItem{
			Id: r.Id,

			HeartRate: r.HeartRate,
			Systolic:  r.Systolic,
			Diastolic: r.Diastolic,
			Spo2:      r.Spo2,
			Posture:   r.Posture,

			RecordTime: r.RecordTime, // RPC 已格式化为 RFC3339 字符串

		})
	}

	return &types.ListRecordsResp{List: list}, nil
}
