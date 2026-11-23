package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"
	"Nurse/app/elderly/model"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateElderlyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateElderlyLogic {
	return &CreateElderlyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateElderlyLogic) CreateElderly(in *elderly.CreateElderlyReq) (*elderly.CreateElderlyResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.GetRealName() == "" {
		return nil, errors.New("real_name is required")
	}
	if g := in.GetGender(); g != 0 && g != 1 && g != 2 {
		return nil, errors.New("gender must be 0(unknown) / 1(male) / 2(female)")
	}
	if cl := in.GetCareLevel(); cl < 1 || cl > 5 {
		return nil, errors.New("care_level must be between 1 and 5")
	}

	// 2. 解析出生日期：仅接受 YYYY-MM-DD
	var birthDate sql.NullTime
	if bs := in.GetBirthDate(); bs != "" {
		t, err := time.ParseInLocation("2006-01-02", bs, time.Local)
		if err != nil {
			return nil, errors.Wrap(err, "invalid birth_date, expected YYYY-MM-DD")
		}
		birthDate = sql.NullTime{Time: t, Valid: true}
	}

	// 3. 慢病史序列化（DB 中存 JSON 数组字符串）
	var chronicDisease sql.NullString
	if list := in.GetChronicDisease(); len(list) > 0 {
		b, err := json.Marshal(list)
		if err != nil {
			return nil, errors.Wrap(err, "marshal chronic_disease failed")
		}
		chronicDisease = sql.NullString{String: string(b), Valid: true}
	}

	// 4. 构造实体
	data := &model.ElderlyArchive{
		RealName:         in.GetRealName(),
		Gender:           int64(in.GetGender()),
		BirthDate:        birthDate,
		IdCardMasked:     sql.NullString{}, // 请求里没有身份证，先留空
		Mobile:           nullStr(in.GetMobile()),
		Address:          nullStr(in.GetAddress()),
		EmergencyContact: nullStr(in.GetEmergencyContact()),
		EmergencyMobile:  nullStr(in.GetEmergencyMobile()),
		ChronicDisease:   chronicDisease,
		CareLevel:        int64(in.GetCareLevel()),
		Status:           1, // 1 正常
		Version:          1, // 乐观锁初始版本
	}

	var newId int64
	err := l.svcCtx.ElderlyArchiveModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		res, err := l.svcCtx.ElderlyArchiveModel.Insert(ctx, session, data)
		if err != nil {
			return err
		}
		newId, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return nil, err
	}

	return &elderly.CreateElderlyResp{
		Id: newId,
	}, nil
}

func nullStr(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
