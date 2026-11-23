// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package training

import (
	"context"

	"Nurse/app/nurse/cmd/api/internal/svc"
	"Nurse/app/nurse/cmd/api/internal/types"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/ctxdata"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTrainingRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 个人考核记录
func NewListTrainingRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTrainingRecordsLogic {
	return &ListTrainingRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTrainingRecordsLogic) ListTrainingRecords() (resp *types.ListTrainingRecordsResp, err error) {
	// todo: add your logic here and delete this line

	// 1. 取当前登录护士 ID（无入参，说明是"我的培训记录"）
	nurseId := ctxdata.GetUidFromCtx(l.ctx)
	if nurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "未获取到登录护士信息")
	}

	// 2. 调用 RPC
	rpcResp, err := l.svcCtx.NurseRpc.ListTrainingRecords(l.ctx, &nurse.ListTrainingRecordsReq{
		NurseId: nurseId,
	})
	if err != nil {
		l.Errorf("调用 RPC ListTrainingRecords 失败: nurseId=%d, err=%v", nurseId, err)
		return nil, xerr.NewErrMsg(err.Error())
	}

	// 3. RPC 响应 → HTTP 类型
	//    TrainingInfo 与 TrainingRecordItem 字段一一对应，直接透传
	list := make([]types.TrainingRecordItem, 0, len(rpcResp.GetList()))
	for _, item := range rpcResp.GetList() {
		list = append(list, types.TrainingRecordItem{
			Id:           item.GetId(),
			CourseName:   item.GetCourseName(),
			ExamMonth:    item.GetExamMonth(),
			Score:        item.GetScore(),
			AttemptCount: item.GetAttemptCount(),
			Passed:       item.GetPassed(),
			CreatedAt:    item.GetCreatedAt(),
		})
	}

	return &types.ListTrainingRecordsResp{List: list}, nil

}
