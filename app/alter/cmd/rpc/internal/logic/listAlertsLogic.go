package logic

import (
	"context"

	"Nurse/app/alter/cmd/rpc/alert"
	"Nurse/app/alter/cmd/rpc/internal/svc"
	"Nurse/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const alertTimeLayout = "2006-01-02 15:04:05"

type ListAlertsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAlertsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAlertsLogic {
	return &ListAlertsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAlertsLogic) ListAlerts(in *alert.ListAlertsReq) (*alert.ListAlertsResp, error) {
	// todo: add your logic here and delete this line

	// 1. 分页参数归一化
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 2. 构建查询条件：handled 0=未处理 1=已处理 其他=全部
	builder := l.svcCtx.AlertEventModel.SelectBuilder()
	if in.Handled == 0 || in.Handled == 1 {
		builder = builder.Where("handled = ?", in.Handled)
	}

	// 3. 分页查询（内部已叠加 del_state = DelStateNo）
	rows, total, err := l.svcCtx.AlertEventModel.FindPageListByPageWithTotal(
		l.ctx, builder, int64(page), int64(pageSize), "id DESC",
	)
	if err != nil {
		l.Errorf("list alerts failed, req=%+v, err=%v", in, err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "query alert list failed: %v", err)
	}

	// 4. 转换为返回结构
	list := make([]*alert.AlertItem, 0, len(rows))

	for _, row := range rows {
		item := &alert.AlertItem{
			Id:        int64(row.Id),
			AlertType: row.AlertType,
			Severity:  int32(row.Severity),
			Handled:   int32(row.Handled),
			// 关键差异：CreatedAt 是 string，需格式化 time.Time
			CreatedAt: row.CreateTime.Format(alertTimeLayout),
		}
		if row.ElderlyId.Valid {
			item.ElderlyId = row.ElderlyId.Int64
		}
		if row.NurseId.Valid {
			item.NurseId = row.NurseId.Int64
		}
		if row.OrderId.Valid {
			item.OrderId = row.OrderId.Int64
		}
		if row.Content.Valid {
			item.Content = row.Content.String
		}
		list = append(list, item)
	}

	return &alert.ListAlertsResp{
		List:  list,
		Total: total,
	}, nil
}
