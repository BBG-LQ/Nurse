package logic

import (
	"context"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/elderly/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListElderlyByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListElderlyByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListElderlyByIdsLogic {
	return &ListElderlyByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const maxBatchIds = 200 // 单次批量查询上限，防止 IN 过大
func (l *ListElderlyByIdsLogic) ListElderlyByIds(in *elderly.ListElderlyByIdsReq) (*elderly.ListElderlyByIdsResp, error) {
	// todo: add your logic here and delete this line

	// 1. 空入参：直接返回空列表，不查库
	if len(in.GetIds()) == 0 {
		return &elderly.ListElderlyByIdsResp{List: []*elderly.ElderlyInfo{}}, nil
	}
	if len(in.GetIds()) > maxBatchIds {
		return nil, status.Errorf(codes.InvalidArgument, "too many ids, max %d", maxBatchIds)
	}

	// 2. 去重 + 过滤非法 id（<=0），同时收集合法 id
	seen := make(map[uint64]struct{}, len(in.GetIds()))
	ids := make([]uint64, 0, len(in.GetIds()))
	for _, id := range in.GetIds() {
		if id <= 0 {
			continue
		}
		u := uint64(id)
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		ids = append(ids, u)
	}
	if len(ids) == 0 {
		return &elderly.ListElderlyByIdsResp{List: []*elderly.ElderlyInfo{}}, nil
	}

	// 3. 构造查询：id IN (?...)，按 id DESC 返回
	builder := l.svcCtx.ElderlyArchiveModel.SelectBuilder().
		Where("id IN (?)", ids)

	rows, err := l.svcCtx.ElderlyArchiveModel.FindAll(l.ctx, builder, "id DESC")
	if err != nil {
		return nil, err
	}

	// 4. 转 proto
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

	return &elderly.ListElderlyByIdsResp{

		List: list,
	}, nil
}
