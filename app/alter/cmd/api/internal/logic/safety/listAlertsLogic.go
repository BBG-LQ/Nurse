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

type ListAlertsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 警报列表
func NewListAlertsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAlertsLogic {
	return &ListAlertsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAlertsLogic) ListAlerts(req *types.ListAlertsReq) (resp *types.ListAlertsResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 分页参数防御（tag 默认值只在显式传参时生效，空 body 时仍需兜底）
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 2. types.Req -> rpc Req 转换
	rpcResp, err := l.svcCtx.AlertRpc.ListAlerts(l.ctx, &alert.ListAlertsReq{
		Handled:  req.Handled,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("call rpc ListAlerts failed, req=%+v, err=%v", req, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "list alerts failed: %v", err)
	}

	// 3. rpc Resp -> types.Resp 转换（切片需逐条拷贝）
	list := make([]types.AlertItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.AlertItem{
			Id:        item.Id,
			ElderlyId: item.ElderlyId,
			NurseId:   item.NurseId,
			OrderId:   item.OrderId,
			AlertType: item.AlertType,
			Severity:  item.Severity,
			Content:   item.Content,
			Handled:   item.Handled,
			CreatedAt: item.CreatedAt,
		})
	}

	return &types.ListAlertsResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
