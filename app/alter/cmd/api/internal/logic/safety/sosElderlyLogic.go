// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package safety

import (
	"context"

	"Nurse/app/alter/cmd/api/internal/svc"
	"Nurse/app/alter/cmd/api/internal/types"
	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SosElderlyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 老人一键呼救
func NewSosElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SosElderlyLogic {
	return &SosElderlyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SosElderlyLogic) SosElderly(req *types.SosElderlyReq) (resp *types.SosElderlyResp, err error) {
	// todo: add your logic here and delete this line

	// 1. API 层基础校验（快速失败，避免无谓 RPC 往返）
	if req.ElderlyId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "老人ID不能为空")
	}

	// 2. types.Req -> rpc Req 转换
	rpcResp, err := l.svcCtx.AlertRpc.SosElderly(l.ctx, &alert.SosElderlyReq{
		ElderlyId: req.ElderlyId,
		Content:   req.Content,
	})
	if err != nil {
		l.Errorf("call rpc SosElderly failed, elderly_id=%d, err=%v", req.ElderlyId, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "sos elderly failed: %v", err)
	}

	// 3. rpc Resp -> types.Resp 转换
	return &types.SosElderlyResp{
		AlertId: rpcResp.AlertId,
	}, nil
}
