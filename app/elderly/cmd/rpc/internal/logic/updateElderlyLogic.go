package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"
	"Nurse/app/elderly/model"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateElderlyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrUserNoExistsError = xerr.NewErrMsg("用户不存在")

func NewUpdateElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateElderlyLogic {
	return &UpdateElderlyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateElderlyLogic) UpdateElderly(in *elderly.UpdateElderlyReq) (*elderly.UpdateElderlyResp, error) {
	// todo: add your logic here and delete this line
	var Ok bool

	// 1. 参数校验
	if cl := in.GetCareLevel(); cl < 1 || cl > 5 {
		return nil, errors.New("care_level must be between 1 and 5")
	}

	// 慢病史序列化（DB 中存 JSON 数组字符串）
	var chronicDisease sql.NullString
	if list := in.GetChronicDisease(); len(list) > 0 {
		b, err := json.Marshal(list)
		if err != nil {
			return nil, errors.Wrap(err, "marshal chronic_disease failed")
		}
		chronicDisease = sql.NullString{String: string(b), Valid: true}
	}

	//find

	m, err := l.svcCtx.ElderlyArchiveModel.FindOne(l.ctx, uint64(in.Id))

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "GetUserInfo find user db err , id:%d , err:%v", in.Id, err)
	}
	if m == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, " id:%d ", in.Id)
	}

	//
	m.Address = nullStr(in.GetAddress())
	m.EmergencyContact = nullStr(in.GetEmergencyContact())
	m.EmergencyMobile = nullStr(in.GetEmergencyMobile())
	m.CareLevel = int64(in.GetCareLevel())
	m.ChronicDisease = chronicDisease

	_, err = l.svcCtx.ElderlyArchiveModel.Update(l.ctx, nil, m)

	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError),
			"UpdateElderly update user db err, id:%d, err:%v", in.Id, err)
	}
	Ok = true

	return &elderly.UpdateElderlyResp{
		Ok: Ok,
	}, nil
}
