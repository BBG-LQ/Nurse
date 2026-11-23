package logic

import (
	"context"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	trainingRecordsLimit = 100
	timeLayout           = "2006-01-02 15:04:05"
)

type ListTrainingRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTrainingRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTrainingRecordsLogic {
	return &ListTrainingRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTrainingRecordsLogic) ListTrainingRecords(in *nurse.ListTrainingRecordsReq) (*nurse.ListTrainingRecordsResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不合法")
	}

	// 2. 按护士ID查询（del_state 由 FindAll 自动追加）
	builder := l.svcCtx.TrainingRecordModel.SelectBuilder().Where(squirrel.Eq{
		"nurse_id": in.NurseId,
	})

	// 3. 按考试月份倒序，同年月按创建时间倒序，最新的在前
	records, err := l.svcCtx.TrainingRecordModel.FindAll(
		l.ctx, builder.Limit(trainingRecordsLimit), "exam_month DESC, id DESC",
	)
	if err != nil {
		l.Errorf("查询培训记录失败: nurseId=%d, err=%v", in.NurseId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 4. 组装返回（字段名请按实际 TrainingInfo 定义调整）
	list := make([]*nurse.TrainingInfo, 0, len(records))

	for _, r := range records {
		info := &nurse.TrainingInfo{
			Id:           int64(r.Id),
			CourseName:   r.CourseName,
			ExamMonth:    r.ExamMonth,
			AttemptCount: int32(r.AttemptCount),
			Passed:       r.Passed == 1,
			CreatedAt:    r.CreateTime.Format(timeLayout),
		}
		// Score 是 sql.NullInt64，需判空后再赋值
		if r.Score.Valid {
			info.Score = int32(r.Score.Int64)
		}
		list = append(list, info)
	}

	return &nurse.ListTrainingRecordsResp{
		List: list,
	}, nil
}
