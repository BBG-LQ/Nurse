package logic

import (
	"context"
	"errors"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/nurse/model"
	"Nurse/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// 乐观锁冲突时的最大重试次数
const incrLoadMaxRetries = 3

type IncrLoadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIncrLoadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncrLoadLogic {
	return &IncrLoadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IncrLoadLogic) IncrLoad(in *nurse.IncrLoadReq) (*nurse.IncrLoadResp, error) {
	// todo: add your logic here and delete this line

	// 1. 参数校验
	if in.NurseId <= 0 {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "护士ID不合法")
	}
	// delta 为 0 视为幂等空操作，直接返回成功
	if in.Delta == 0 {
		return &nurse.IncrLoadResp{Ok: true}, nil
	}

	// 2. 乐观锁重试循环
	for i := 0; i < incrLoadMaxRetries; i++ {
		profile, err := l.svcCtx.NurseProfileModel.FindOne(l.ctx, uint64(in.NurseId))
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil, xerr.NewErrCodeMsg(xerr.DataNoExistError, "护士不存在")
			}
			l.Errorf("查询护士信息失败: nurseId=%d, err=%v", in.NurseId, err)
			return nil, xerr.NewErrCode(xerr.DbError)
		}

		// 3. 计算新负载，禁止为负
		newLoad := profile.CurrentLoad + int64(in.Delta)
		if newLoad < 0 {
			l.Errorf("负载变更会导致负数: nurseId=%d, currentLoad=%d, delta=%d",
				in.NurseId, profile.CurrentLoad, in.Delta)
			return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "负载不足，无法扣减")
		}
		profile.CurrentLoad = newLoad

		// 4. 乐观锁更新
		err = l.svcCtx.NurseProfileModel.UpdateWithVersion(l.ctx, nil, profile)
		if err == nil {
			return &nurse.IncrLoadResp{Ok: true}, nil
		}
		if errors.Is(err, model.ErrNoRowsUpdate) {
			// 版本冲突，重新读取后重试
			l.Infof("更新负载版本冲突，重试: nurseId=%d, attempt=%d", in.NurseId, i+1)
			continue
		}
		l.Errorf("更新护士负载失败: nurseId=%d, delta=%d, err=%v", in.NurseId, in.Delta, err)
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 重试耗尽
	return nil, xerr.NewErrCodeMsg(xerr.DbUpdateAffectedZeroError, "更新负载失败，请重试")
}
