package logic

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/nurse/model"
	"Nurse/pkg/globalkey"
	"Nurse/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

// 及格线（可根据业务调整）
const passScore int32 = 60

// 考试月份格式：YYYY-MM
var examMonthPattern = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)

type SubmitTrainingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTrainingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTrainingLogic {
	return &SubmitTrainingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitTrainingLogic) SubmitTraining(in *nurse.SubmitTrainingReq) (*nurse.SubmitTrainingResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不合法")
	}
	if in.CourseName == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "课程名称不能为空")
	}
	if !examMonthPattern.MatchString(in.ExamMonth) {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "考试月份格式应为 YYYY-MM")
	}
	if in.Score < 0 || in.Score > 100 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "成绩必须在 0-100 之间")
	}

	// 2. 校验护士是否存在
	if _, err := l.svcCtx.NurseProfileModel.FindOne(l.ctx, uint64(in.NurseId)); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "护士不存在")
		}
		l.Errorf("查询护士信息失败: nurseId=%d, err=%v", in.NurseId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	passed := in.Score >= passScore
	passedInt := int64(0)
	if passed {
		passedInt = 1
	}

	// 3. 查询该护士、该课程、该月份的培训记录
	builder := l.svcCtx.TrainingRecordModel.SelectBuilder().Where(squirrel.Eq{
		"nurse_id":    in.NurseId,
		"course_name": in.CourseName,
		"exam_month":  in.ExamMonth,
	})
	records, err := l.svcCtx.TrainingRecordModel.FindAll(l.ctx, builder, "")
	if err != nil {
		l.Errorf("查询培训记录失败: err=%v", err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 4. 不存在则新增
	if len(records) == 0 {
		newRecord := &model.TrainingRecord{
			NurseId:      uint64(in.NurseId),
			CourseName:   in.CourseName,
			ExamMonth:    in.ExamMonth,
			Score:        sql.NullInt64{Int64: int64(in.Score), Valid: true},
			AttemptCount: 1,
			Passed:       passedInt,
			Status:       1,
			DeleteTime:   time.Unix(0, 0),
			DelState:     globalkey.DelStateNo,
			Version:      1,
		}
		if _, err := l.svcCtx.TrainingRecordModel.Insert(l.ctx, nil, newRecord); err != nil {
			l.Errorf("新增培训记录失败: err=%v", err)
			return nil, xerr.NewErrCode(xerr.DbError)
		}
		return &nurse.SubmitTrainingResp{
			Passed:       passed,
			AttemptCount: 1,
		}, nil
	}

	// 5. 已通过则直接返回（不再重复考核）
	record := records[0]
	if record.Passed == 1 {
		return &nurse.SubmitTrainingResp{
			Passed:       true,
			AttemptCount: int32(record.AttemptCount),
		}, nil
	}

	// 6. 未通过则累加尝试次数并更新成绩
	record.Score = sql.NullInt64{Int64: int64(in.Score), Valid: true}
	record.AttemptCount += 1
	record.Passed = passedInt
	if err := l.svcCtx.TrainingRecordModel.UpdateWithVersion(l.ctx, nil, record); err != nil {
		if errors.Is(err, model.ErrNoRowsUpdate) {
			return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "记录已被修改，请重试")
		}
		l.Errorf("更新培训记录失败: err=%v", err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	return &nurse.SubmitTrainingResp{
		Passed:       passed,
		AttemptCount: int32(record.AttemptCount),
	}, nil
}
