// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package training

import (
	"context"
	"strings"

	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitTrainingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 提交考核
func NewSubmitTrainingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTrainingLogic {
	return &SubmitTrainingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitTrainingLogic) SubmitTraining(req *types.SubmitTrainingReq) (resp *types.SubmitTrainingResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 取当前登录护士 ID（types 里没有 nurseId，说明由 JWT 注入）
	nurseId := ctxdata.GetUidFromCtx(l.ctx)
	if nurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "未获取到登录护士信息")
	}

	// 2. 基础参数校验（快速失败，减少无意义 RPC 调用）
	courseName := strings.TrimSpace(req.CourseName)
	if courseName == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "课程名称不能为空")
	}
	examMonth := strings.TrimSpace(req.ExamMonth)
	if examMonth == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "考试月份不能为空")
	}

	// 3. 构造 RPC 请求并调用
	rpcResp, err := l.svcCtx.NurseRpc.SubmitTraining(l.ctx, &nurse.SubmitTrainingReq{
		NurseId:    nurseId,
		CourseName: courseName,
		ExamMonth:  examMonth,
		Score:      req.Score,
	})
	if err != nil {
		l.Errorf("调用 RPC SubmitTraining 失败: nurseId=%d, req=%+v, err=%v", nurseId, req, err)
		return nil, xerr.NewErrMsg(err.Error())
	}

	// 4. RPC 响应 → HTTP 类型
	return &types.SubmitTrainingResp{
		Passed:       rpcResp.GetPassed(),
		AttemptCount: rpcResp.GetAttemptCount(),
	}, nil
}
