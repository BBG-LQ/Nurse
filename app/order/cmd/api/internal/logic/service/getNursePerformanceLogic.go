// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package service

import (
	"context"

	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNursePerformanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 护士绩效汇总
func NewGetNursePerformanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNursePerformanceLogic {
	return &GetNursePerformanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNursePerformanceLogic) GetNursePerformance(req *types.PerformanceReq) (resp *types.PerformanceResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	//    nurseId 是 query 参数（form tag），缺省时为 0，必须拦截
	if req.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不合法")
	}

	// 2. 越权校验：护士只能查看自己的绩效
	//    uid=0（管理端/未鉴权）时放行；与 GetOrder 的宽松策略保持一致
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.TokenExpireError, "未登录或登录已过期")
	}

	if uid > 0 && uid != req.NurseId {
		l.Infof("护士尝试查看他人绩效: reqNurseId=%d, reqUid=%d", req.NurseId, uid)
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "无权查看该护士绩效")
	}

	// 3. 调用 RPC
	rpcResp, err := l.svcCtx.OrderRpc.GetNursePerformance(l.ctx, &order.GetNursePerformanceReq{
		NurseId: req.NurseId,
	})
	if err != nil {
		l.Errorf("调用 GetNursePerformance RPC 失败: nurseId=%d, err=%v", req.NurseId, err)
		return nil, err
	}

	// 4. 字段映射（一一对应）
	return &types.PerformanceResp{
		OrderCount:     rpcResp.OrderCount,
		ServiceMinutes: rpcResp.ServiceMinutes,
		ScoreAvg:       rpcResp.ScoreAvg,
	}, nil
}
