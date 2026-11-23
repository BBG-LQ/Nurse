package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"strings"
	"time"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"
	"Nurse/app/elderly/model"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateMedicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMedicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMedicationLogic {
	return &CreateMedicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var ErrRemind_timesError = xerr.NewErrMsg("marshal remind_times failed")

func (l *CreateMedicationLogic) CreateMedication(in *elderly.CreateMedicationReq) (*elderly.CreateMedicationResp, error) {
	// todo: add your logic here and delete this line

	// 1. 必填校验
	if in.GetElderlyId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "elderly_id is required")
	}
	if in.GetNurseId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nurse_id is required")
	}
	if strings.TrimSpace(in.GetDrugName()) == "" {
		return nil, status.Error(codes.InvalidArgument, "drug_name is required")
	}
	if strings.TrimSpace(in.GetDosage()) == "" {
		return nil, status.Error(codes.InvalidArgument, "dosage is required")
	}
	if strings.TrimSpace(in.GetFrequency()) == "" {
		return nil, status.Error(codes.InvalidArgument, "frequency is required")
	}

	// 2. 校验老人是否存在
	if _, err := l.svcCtx.ElderlyArchiveModel.FindOne(l.ctx, uint64(in.GetElderlyId())); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "elderly not found")
		}
		return nil, err
	}

	// 3. 解析起止日期（YYYY-MM-DD）
	startDate, err := parseNullDate(in.GetStartDate(), "start_date")
	if err != nil {
		return nil, err
	}
	endDate, err := parseNullDate(in.GetEndDate(), "end_date")
	if err != nil {
		return nil, err
	}
	if startDate.Valid && endDate.Valid && endDate.Time.Before(startDate.Time) {
		return nil, status.Error(codes.InvalidArgument, "end_date must be >= start_date")
	}

	// 4. remind_times 序列化为 JSON，并校验 HH:MM 格式
	remindTimes, err := buildRemindTimes(in.GetRemindTimes())
	if err != nil {
		return nil, err
	}

	// 5. 构造实体
	data := &model.ElderlyMedication{
		ElderlyId:   uint64(in.GetElderlyId()),
		NurseId:     uint64(in.GetNurseId()),
		DrugName:    strings.TrimSpace(in.GetDrugName()),
		Dosage:      nullStr(strings.TrimSpace(in.GetDosage())),
		Frequency:   nullStr(strings.TrimSpace(in.GetFrequency())),
		RemindTimes: remindTimes,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      1, // 1 用药中
		Version:     1,
	}

	// 6. 事务插入
	var newId int64
	err = l.svcCtx.ElderlyMedicationModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		res, err := l.svcCtx.ElderlyMedicationModel.Insert(ctx, session, data)
		if err != nil {
			return err
		}
		newId, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return nil, err
	}

	return &elderly.CreateMedicationResp{
		Id: newId,
	}, nil
}

func parseNullDate(s, field string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{}, nil
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return sql.NullTime{}, status.Errorf(codes.InvalidArgument, "%s must be YYYY-MM-DD", field)
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}

func buildRemindTimes(list []string) (sql.NullString, error) {
	if len(list) == 0 {
		return sql.NullString{}, nil
	}
	clean := make([]string, 0, len(list))
	for _, t := range list {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, err := time.Parse("15:04", t); err != nil {
			return sql.NullString{}, status.Errorf(codes.InvalidArgument, "invalid remind_time %q, expected HH:MM", t)
		}
		clean = append(clean, t)
	}
	if len(clean) == 0 {
		return sql.NullString{}, nil
	}
	b, err := json.Marshal(clean)
	if err != nil {
		return sql.NullString{}, errors.Wrapf(ErrRemind_timesError, "marshal remind_times failed +%v", clean)
	}
	return sql.NullString{String: string(b), Valid: true}, nil
}
func parseRemindTimes(ns sql.NullString) []string {
	// NullString.Valid 为 false 表示 DB 里是 NULL；String 为空串也视为“没有提醒时间”
	// 两种情况统一返回空切片（而非 nil），方便前端统一按数组处理
	if !ns.Valid || ns.String == "" {
		return []string{}
	}

	// 尝试按 JSON 数组反序列化（正常写入路径的数据都是这种格式）
	var list []string
	if err := json.Unmarshal([]byte(ns.String), &list); err != nil {
		// JSON 解析失败：兼容老数据 —— 若内容不是以 '[' 开头，按逗号分隔再拆一次
		if s := strings.TrimSpace(ns.String); s != "" && !strings.HasPrefix(s, "[") {
			parts := strings.Split(s, ",")
			out := make([]string, 0, len(parts))
			for _, p := range parts {
				if p = strings.TrimSpace(p); p != "" {
					out = append(out, p)
				}
			}
			return out
		}
		// 既不是合法 JSON、也不是逗号分隔串 → 脏数据，打日志后降级为空
		logx.Errorf("parse remind_times failed, raw=%q, err=%v", ns.String, err)
		return []string{}
	}

	// JSON 是 "null" 或 "[]" 时 list 可能为 nil → 统一成空切片
	if list == nil {
		return []string{}
	}
	return list
}
