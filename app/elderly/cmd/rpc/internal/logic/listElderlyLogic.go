package logic

import (
	"context"
	"strings"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListElderlyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListElderlyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListElderlyLogic {
	return &ListElderlyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

func (l *ListElderlyLogic) ListElderly(in *elderly.ListElderlyReq) (*elderly.ListElderlyResp, error) {
	// todo: add your logic here and delete this line

	page := in.GetPage()
	if page < 1 {
		page = defaultPage
	}
	pageSize := in.GetPageSize()
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	builder := l.svcCtx.ElderlyArchiveModel.SelectBuilder()

	// 关键字模糊搜索：姓名 / 手机号
	if kw := strings.TrimSpace(in.GetKeyword()); kw != "" {
		like := "%" + kw + "%"
		builder = builder.Where("(real_name LIKE ? OR mobile LIKE ?)", like, like)
	}

	rows, total, err := l.svcCtx.ElderlyArchiveModel.FindPageListByPageWithTotal(
		l.ctx, builder, int64(page), int64(pageSize), "id DESC",
	)
	if err != nil {
		return nil, err
	}

	list := make([]*elderly.ElderlyInfo, 0, len(rows))
	for _, m := range rows {
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

		list = append(list, info)
	}

	return &elderly.ListElderlyResp{
		List:  list,
		Total: total,
	}, nil
}
