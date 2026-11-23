// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package training

import (
	"context"
	"strings"

	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUnqualifiedNursesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 三次不合格名单
func NewListUnqualifiedNursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUnqualifiedNursesLogic {
	return &ListUnqualifiedNursesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUnqualifiedNursesLogic) ListUnqualifiedNurses(req *types.UnqualifiedReq) (resp *types.UnqualifiedResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	examMonth := strings.TrimSpace(req.ExamMonth)
	if examMonth == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "考试月份不能为空")
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.NurseRpc.ListUnqualifiedNurses(l.ctx, &nurse.ListUnqualifiedNursesReq{
		ExamMonth: examMonth,
	})
	if err != nil {
		l.Errorf("调用 RPC ListUnqualifiedNurses 失败: examMonth=%s, err=%v", examMonth, err)
		return nil, xerr.NewErrMsg(err.Error())
	}

	// 3. RPC 响应 → HTTP 类型
	//    proto 里 nurse_ids 是 []int64，与 types.NurseIds 完全一致，直接透传
	nurseIds := rpcResp.GetNurseIds()
	if nurseIds == nil {
		nurseIds = []int64{}
	}

	return &types.UnqualifiedResp{
		NurseIds: nurseIds,
	}, nil
}
