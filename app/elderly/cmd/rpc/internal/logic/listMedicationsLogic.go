package logic

import (
	"context"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMedicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMedicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMedicationsLogic {
	return &ListMedicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMedicationsLogic) ListMedications(in *elderly.ListMedicationsReq) (*elderly.ListMedicationsResp, error) {
	// todo: add your logic here and delete this line

	builder := l.svcCtx.ElderlyMedicationModel.SelectBuilder().
		Where("elderly_id = ?", uint64(in.GetElderlyId()))

	rows, err := l.svcCtx.ElderlyMedicationModel.FindAll(
		l.ctx, builder, "start_date DESC, id DESC",
	)
	if err != nil {
		return nil, err
	}

	list := make([]*elderly.MedicationInfo, 0, len(rows))

	for _, m := range rows {

		info := &elderly.MedicationInfo{
			Id:          int64(m.Id),
			ElderlyId:   int64(m.ElderlyId),
			NurseId:     int64(m.NurseId),
			DrugName:    m.DrugName,
			Dosage:      m.Dosage.String,    // 无效时为空串
			Frequency:   m.Frequency.String, // 无效时为空串
			RemindTimes: parseRemindTimes(m.RemindTimes),
			Status:      int32(m.Status),
		}
		list = append(list, info)
	}

	return &elderly.ListMedicationsResp{
		List: list,
	}, nil
}
