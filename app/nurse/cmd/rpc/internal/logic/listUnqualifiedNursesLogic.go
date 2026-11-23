package logic

import (
	"context"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

const unqualifiedRecordsLimit = 5000

type ListUnqualifiedNursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUnqualifiedNursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUnqualifiedNursesLogic {
	return &ListUnqualifiedNursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUnqualifiedNursesLogic) ListUnqualifiedNurses(in *nurse.ListUnqualifiedNursesReq) (*nurse.ListUnqualifiedNursesResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验：考试月份必填
	if in.ExamMonth == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "考试月份不能为空")
	}

	// 2. 查询该月份未通过的培训记录（passed = 0，del_state 由 FindAll 自动追加）
	builder := l.svcCtx.TrainingRecordModel.SelectBuilder().Where(squirrel.Eq{
		"exam_month": in.ExamMonth,
		"passed":     0,
	})

	records, err := l.svcCtx.TrainingRecordModel.FindAll(
		l.ctx, builder.Limit(unqualifiedRecordsLimit), "nurse_id ASC",
	)
	if err != nil {
		l.Errorf("查询未通过培训记录失败: examMonth=%s, err=%v", in.ExamMonth, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 3. 去重：同一护士该月份可能有多门课程未通过，只返回一次
	seen := make(map[uint64]struct{}, len(records))

	nurseIds := make([]int64, 0, len(records))

	for _, r := range records {
		if _, ok := seen[r.NurseId]; ok {
			continue
		}
		seen[r.NurseId] = struct{}{}
		nurseIds = append(nurseIds, int64(r.NurseId))
	}

	return &nurse.ListUnqualifiedNursesResp{
		NurseIds: nurseIds,
	}, nil
}
