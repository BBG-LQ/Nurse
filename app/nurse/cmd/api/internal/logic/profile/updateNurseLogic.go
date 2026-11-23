// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"context"
	"strings"

	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNurseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新档案/触发复核
func NewUpdateNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNurseLogic {
	return &UpdateNurseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNurseLogic) UpdateNurse(req *types.UpdateNurseReq) (resp *types.UpdateNurseResp, err error) {
	// todo: add your logic here and delete this line
	// 1. path 参数基础校验
	if req.Id <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不合法")
	}

	// 2. 构造 RPC 请求
	//    注意：types 里用零值表示"未传"，这里只对显式传入的字段做 trim / 透传，
	//    真正的"字段是否要更新"的判定交给 RPC 层（见下方说明）
	rpcReq := &nurse.UpdateNurseReq{
		Id:            req.Id,
		SkillTags:     req.SkillTags,
		ServiceArea:   strings.TrimSpace(req.ServiceArea),
		Status:        req.Status,
		LicenseExpire: strings.TrimSpace(req.LicenseExpire),
	}

	// 3. 调用 RPC
	rpcResp, err := l.svcCtx.NurseRpc.UpdateNurse(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("调用 RPC UpdateNurse 失败: req=%+v, err=%v", req, err)
		return nil, xerr.NewErrMsg(err.Error())
	}

	// 4. RPC 响应 → HTTP 类型
	return &types.UpdateNurseResp{
		Ok: rpcResp.GetOk(),
	}, nil
}
