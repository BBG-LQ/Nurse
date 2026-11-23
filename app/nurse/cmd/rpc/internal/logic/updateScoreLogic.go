package logic

import (
	"context"
	"errors"
	"math"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/nurse/model"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// 评分取值范围，按业务调整（此处假设 0~5 分制）
const (
	minScoreAvg = 0.0
	maxScoreAvg = 5.0
)

// 乐观锁冲突时的最大重试次数
const updateScoreMaxRetries = 3

type UpdateScoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateScoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateScoreLogic {
	return &UpdateScoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateScoreLogic) UpdateScore(in *nurse.UpdateScoreReq) (*nurse.UpdateScoreResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不合法")
	}
	// 排除 NaN / +Inf / -Inf，这类值写入 DB 会污染排序与聚合
	if math.IsNaN(in.ScoreAvg) || math.IsInf(in.ScoreAvg, 0) {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "评分值不合法")
	}
	if in.ScoreAvg < minScoreAvg || in.ScoreAvg > maxScoreAvg {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "评分值超出允许范围")
	}

	// 2. 乐观锁重试循环
	for i := 0; i < updateScoreMaxRetries; i++ {
		profile, err := l.svcCtx.NurseProfileModel.FindOne(l.ctx, uint64(in.NurseId))
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "护士不存在")
			}
			l.Errorf("查询护士信息失败: nurseId=%d, err=%v", in.NurseId, err)
			return nil, xerr.NewErrCode(xerr.DbError)
		}

		// 3. 幂等：评分未变化直接返回成功，省掉一次写库 + 缓存清理
		if profile.ScoreAvg == in.ScoreAvg {
			return &nurse.UpdateScoreResp{Ok: true}, nil
		}
		profile.ScoreAvg = in.ScoreAvg

		// 4. 乐观锁更新
		err = l.svcCtx.NurseProfileModel.UpdateWithVersion(l.ctx, nil, profile)
		if err == nil {
			return &nurse.UpdateScoreResp{Ok: true}, nil
		}
		if errors.Is(err, model.ErrNoRowsUpdate) {
			l.Infof("更新评分版本冲突，重试: nurseId=%d, attempt=%d", in.NurseId, i+1)
			continue
		}
		l.Errorf("更新护士评分失败: nurseId=%d, scoreAvg=%f, err=%v", in.NurseId, in.ScoreAvg, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 重试耗尽
	return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "更新评分失败，请重试")
}
