// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package service

import (
	"context"
	"unicode/utf8"

	"Nurse/app/order/cmd/api/internal/svc"
	"Nurse/app/order/cmd/api/internal/types"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	satisfactionScoreMin   = 1
	satisfactionScoreMax   = 5
	satisfactionCommentMax = 500 // 评论最大字符数，与 RPC 层保持一致
)

type SubmitSatisfactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 满意度评分
func NewSubmitSatisfactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitSatisfactionLogic {
	return &SubmitSatisfactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitSatisfactionLogic) SubmitSatisfaction(req *types.SatisfactionReq) (resp *types.SatisfactionResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验（快速失败，避免无谓 RPC 往返）
	//    go-zero 已把 path:"id" 解析到 req.OrderId
	if req.OrderId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "订单ID不合法")
	}
	if req.Score < satisfactionScoreMin || req.Score > satisfactionScoreMax {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "评分必须为 1-5")
	}
	if utf8.RuneCountInString(req.Comment) > satisfactionCommentMax {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "评论内容过长")
	}

	// 2. 调用 RPC
	//    NurseId 不传（填 0）：以订单实际护士为准，防止客户端伪造；
	//    归属校验（uid == order.ElderlyId）由 RPC 层完成
	rpcResp, err := l.svcCtx.OrderRpc.SubmitSatisfaction(l.ctx, &order.SubmitSatisfactionReq{
		OrderId: req.OrderId,

		Score:   req.Score,
		Comment: req.Comment,
	})
	if err != nil {
		// 不重造错误，交给全局 ErrorHandler 统一响应
		l.Errorf("调用 SubmitSatisfaction RPC 失败: orderId=%d, score=%d, err=%v",
			req.OrderId, req.Score, err)
		return nil, err
	}

	// 3. 映射响应
	return &types.SatisfactionResp{
		Ok: rpcResp.Ok,
	}, nil

}
