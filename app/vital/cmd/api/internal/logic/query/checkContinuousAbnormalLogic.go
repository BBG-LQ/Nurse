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

type CheckContinuousAbnormalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 连续异常判定
func NewCheckContinuousAbnormalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckContinuousAbnormalLogic {
	return &CheckContinuousAbnormalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckContinuousAbnormalLogic) CheckContinuousAbnormal(req *types.CheckAbnormalReq) (resp *types.CheckAbnormalResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if req.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "elderlyId 不能为空")
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.VitalRpc.CheckContinuousAbnormal(l.ctx, &vital.CheckContinuousAbnormalReq{
		ElderlyId: req.ElderlyId,
	})
	if err != nil {
		l.Errorf("rpc CheckContinuousAbnormal failed, elderlyId=%d, err=%v", req.ElderlyId, err)
		return nil, xerr.NewErrCode(xerr.ServerCommonError)
	}

	// 3. 转换为 API 响应
	return &types.CheckAbnormalResp{
		Continuous: rpcResp.Continuous,
		Days:       rpcResp.Days,
	}, nil
}
