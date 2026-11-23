// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"context"

	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNurseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 护士详情
func NewGetNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNurseLogic {
	return &GetNurseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNurseLogic) GetNurse() (resp *types.NurseDetail, err error) {
	// todo: add your logic here and delete this line

	// 1. 取当前登录护士 ID（从 JWT ctx 中取，具体函数名按项目 jwt 中间件实现）
	nurseId := ctxdata.GetUidFromCtx(l.ctx)
	if nurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "未获取到登录护士信息")
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.NurseRpc.GetNurse(l.ctx, &nurse.GetNurseReq{
		Id: nurseId,
	})
	if err != nil {
		l.Errorf("调用 RPC GetNurse 失败: nurseId=%d, err=%v", nurseId, err)
		return nil, xerr.NewErrMsg(err.Error())
	}

	// 3. RPC 响应 → HTTP 类型
	return &types.NurseDetail{
		Id:            rpcResp.Info.Id,
		RealName:      rpcResp.Info.RealName,
		Mobile:        rpcResp.Info.Mobile,
		LicenseNo:     rpcResp.Info.LicenseNo,
		LicenseExpire: rpcResp.Info.LicenseExpire,
		SkillTags:     rpcResp.Info.SkillTags,
		ServiceArea:   rpcResp.Info.ServiceArea,
		CurrentLoad:   rpcResp.Info.CurrentLoad,
		ScoreAvg:      rpcResp.Info.ScoreAvg,
		Status:        rpcResp.Info.Status,
	}, nil
}
