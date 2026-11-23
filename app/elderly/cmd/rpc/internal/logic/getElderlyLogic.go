package logic

import (
	"context"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetElderlyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetElderlyLogic {
	return &GetElderlyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetElderlyLogic) GetElderly(in *elderly.GetElderlyReq) (*elderly.GetElderlyResp, error) {
	// todo: add your logic here and delete this line

	//check
	if in.Id <= 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("id no exists"), "id no exists , in:%+v", in.Id)

	}

	//get info
	m, err := l.svcCtx.ElderlyArchiveModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, err
	}

	info := &elderly.ElderlyInfo{
		Id:       int64(m.Id),
		RealName: m.RealName,
		Gender:   int32(m.Gender),
		// IdCardMasked:     m.IdCardMasked.String,
		Mobile:           m.Mobile.String,
		Address:          m.Address.String,
		EmergencyContact: m.EmergencyContact.String,
		EmergencyMobile:  m.EmergencyMobile.String,
		CareLevel:        int32(m.CareLevel),
		Status:           int32(m.Status),
		// Version:          m.Version,
	}
	return &elderly.GetElderlyResp{
		Info: info,
	}, nil
}
