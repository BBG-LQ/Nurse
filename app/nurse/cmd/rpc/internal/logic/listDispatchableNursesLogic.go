package logic

import (
	"context"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/pkg/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

// 单次返回上限，防止一次拉取过多数据
const dispatchableNursesLimit = 50

type ListDispatchableNursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDispatchableNursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDispatchableNursesLogic {
	return &ListDispatchableNursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDispatchableNursesLogic) ListDispatchableNurses(in *nurse.ListDispatchableNursesReq) (*nurse.ListDispatchableNursesResp, error) {
	// todo: add your logic here and delete this line

	// 1. 基础条件：只取可接单状态（status=1），del_state 由 FindAll 自动追加
	builder := l.svcCtx.NurseProfileModel.SelectBuilder().Where(squirrel.Eq{
		"status": 1,
	})

	// 2. 负载上限：只取当前负载小于 maxLoad 的护士
	if in.MaxLoad > 0 {
		builder = builder.Where("current_load < ?", in.MaxLoad)
	}

	// 3. 服务区域匹配：service_area 为逗号分隔的虚拟区域串，用 LIKE 做包含匹配
	if in.Area != "" {
		builder = builder.Where("service_area LIKE ?", "%"+in.Area+"%")
	}

	// 4. 技能匹配：skill_tags 为逗号分隔的标签串，用 LIKE 做包含匹配
	if in.ServiceType != "" {
		builder = builder.Where("skill_tags LIKE ?", "%"+in.ServiceType+"%")
	}

	// 5. 调度优先级：负载低的优先，负载相同按评分高优先
	records, err := l.svcCtx.NurseProfileModel.FindAll(
		l.ctx, builder.Limit(dispatchableNursesLimit), "current_load ASC, score_avg DESC",
	)
	if err != nil {
		l.Errorf("查询可派单护士失败: serviceType=%s, area=%s, maxLoad=%d, err=%v",
			in.ServiceType, in.Area, in.MaxLoad, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 6. 组装返回（字段名请按实际 NurseBrief 定义调整）
	nurses := make([]*nurse.NurseBrief, 0, len(records))
	for _, r := range records {
		nurses = append(nurses, &nurse.NurseBrief{
			Id:       int64(r.Id),
			RealName: r.RealName,

			ScoreAvg:    r.ScoreAvg,
			CurrentLoad: int32(r.CurrentLoad),
		})
	}

	return &nurse.ListDispatchableNursesResp{
		Nurses: nurses,
	}, nil
}
