// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adverse

import (
	"context"
	"time"

	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AdverseStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 统计图表数据
func NewAdverseStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdverseStatisticsLogic {
	return &AdverseStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdverseStatisticsLogic) AdverseStatistics(req *types.StatisticsReq) (resp *types.StatisticsResp, err error) {
	// todo: add your logic here and delete this line
	// 1. 基础格式校验（非法月份无需浪费 RPC 往返）
	//    空串表示"默认当月"，交给 RPC 层处理，此处放行
	if req.Month != "" && !isValidMonth(req.Month) {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "月份格式错误，应为 YYYY-MM")
	}

	// 2. types.Req -> rpc Req 转换
	rpcResp, err := l.svcCtx.AlertRpc.AdverseStatistics(l.ctx, &alert.AdverseStatisticsReq{
		Month: req.Month,
	})
	if err != nil {
		l.Errorf("call rpc AdverseStatistics failed, month=%q, err=%v", req.Month, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "adverse statistics failed: %v", err)
	}

	// 3. rpc Resp -> types.Resp 转换
	//    两侧 TypeCount 均为 map[string]int32（同类型），可直接赋值
	return &types.StatisticsResp{
		TypeCount: rpcResp.TypeCount,
		Total:     rpcResp.Total,
	}, nil
}

// isValidMonth 校验 YYYY-MM 格式且月分为合法值
func isValidMonth(s string) bool {
	if len(s) != 7 || s[4] != '-' {
		return false
	}
	_, err := time.Parse("2006-01", s)
	return err == nil
}
