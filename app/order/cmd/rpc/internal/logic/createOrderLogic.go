package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"Nurse/app/elderly/cmd/rpc/elderly"
	"Nurse/app/order/cmd/rpc/internal/svc"
	"Nurse/app/order/cmd/rpc/order"
	"Nurse/app/order/model"
	"Nurse/pkg/globalkey"
	"Nurse/pkg/tool"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// 紧急程度
const (
	urgencyNormal = 0
	urgencyUrgent = 1
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验

	scheduledAt, err := parseScheduledAt(in.ScheduledAt)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError,
			"预约时间格式错误，应为 2006-01-02 15:04:05 或 RFC3339")
	}

	// 2. 校验老人是否存在
	if _, err := l.svcCtx.ElderlyRpc.GetElderly(l.ctx, &elderly.GetElderlyReq{
		Id: in.ElderlyId,
	}); err != nil {
		l.Errorf("查询老人信息失败: elderlyId=%d, err=%v", in.ElderlyId, err)
		return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "老人不存在")
	}

	// 3. 组装订单数据
	orderNo := genOrderNo()
	data := &model.ServiceOrder{
		OrderNo:     orderNo,
		ElderlyId:   uint64(in.ElderlyId),
		NurseId:     sql.NullInt64{}, // 待派，暂未分配护士
		ServiceType: in.ServiceType,
		Status:      orderStatusPending,
		Urgency:     int64(in.Urgency),
		ScheduledAt: scheduledAt,
		StartedAt:   sql.NullTime{},
		FinishedAt:  sql.NullTime{},
		Amount:      sql.NullFloat64{Float64: in.Amount, Valid: true},
		DeleteTime:  time.Unix(0, 0),
		DelState:    globalkey.DelStateNo,
		Version:     0,
	}

	// 4. 单表插入，无需事务
	result, err := l.svcCtx.ServiceOrderModel.Insert(l.ctx, nil, data)
	if err != nil {
		l.Errorf("创建订单失败: orderNo=%s, elderlyId=%d, err=%v",
			orderNo, in.ElderlyId, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	id, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取订单自增ID失败: orderNo=%s, err=%v", orderNo, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	return &order.CreateOrderResp{
		OrderId: id,
		OrderNo: orderNo,
	}, nil
}

// parseScheduledAt 解析预约时间，支持空串、MySQL 常用格式、RFC3339
func parseScheduledAt(s string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{}, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return sql.NullTime{Time: t, Valid: true}, nil
		}
	}
	return sql.NullTime{}, errors.New("invalid scheduled_at format")
}

// genOrderNo 生成订单号：SO + 秒级时间戳 + 6位随机数字
// 同秒冲突概率极低，且表上 order_no 有唯一索引兜底
func genOrderNo() string {
	return "SO" + time.Now().Format("20060102150405") + tool.Krand(6, tool.KC_RAND_KIND_NUM)
}
