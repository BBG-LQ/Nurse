package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateNurseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNurseLogic {
	return &UpdateNurseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateNurseLogic) UpdateNurse(in *nurse.UpdateNurseReq) (*nurse.UpdateNurseResp, error) {
	// todo: add your logic here and delete this line

	// ========== 1. 基础参数校验 ==========

	// 1.1 ID 必须为正整数
	if in.GetId() <= 0 {
		return nil, errors.New("id 不合法")
	}

	// 1.2 至少有一个字段需要更新，避免无意义的 DB 写操作
	if len(in.GetSkillTags()) == 0 &&
		strings.TrimSpace(in.GetServiceArea()) == "" &&
		in.GetStatus() == 0 &&
		strings.TrimSpace(in.GetLicenseExpire()) == "" {
		return nil, errors.New("没有需要更新的字段")
	}

	// ========== 2. 查询原记录（确认存在 + 拿到完整数据用于 Update） ==========

	// 2.1 先查出原记录；go-zero 的 Update 是全量更新，因此必须先查出完整行再改字段
	result, err := l.svcCtx.NurseProfileModel.FindOne(l.ctx, uint64(in.GetId()))
	if err != nil {
		// 2.2 记录不存在
		if errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("护士不存在")
		}
		// 2.3 其它 DB 错误：记日志、返回通用提示
		l.Errorf("find nurse profile failed, id=%d, err=%v", in.GetId(), err)
		return nil, errors.New("系统繁忙，请稍后重试")
	}
	if result == nil {
		return nil, errors.New("护士不存在")
	}

	// ========== 3. 逐字段处理（有传才更新） ==========
	// 注意：proto3 未启用 optional 时，标量字段无法区分"未传"与"传了零值"，
	//      这里按业务约定：0 / "" / 空数组 视为"不更新"，需要清空请单独设计接口。

	// 3.1 技能标签：清洗后序列化为 JSON 存入 sql.NullString
	if len(in.GetSkillTags()) > 0 {
		// 3.1.1 数量上限保护
		if len(in.GetSkillTags()) > 20 {
			return nil, errors.New("skill_tags 数量不能超过 20 个")
		}
		// 3.1.2 去空、去首尾空白
		tags := make([]string, 0, len(in.GetSkillTags()))
		for _, tag := range in.GetSkillTags() {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			tags = append(tags, tag)
		}
		// 3.1.3 如果清洗后为空，视为非法（传了但全是空串）
		if len(tags) == 0 {
			return nil, errors.New("skill_tags 内容不能为空")
		}
		// 3.1.4 序列化 JSON
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			l.Errorf("marshal skill_tags failed, id=%d, err=%v", in.GetId(), err)
			return nil, errors.New("系统繁忙，请稍后重试")
		}
		// 3.1.5 赋值：sql.NullString{String, Valid=true} 表示非 NULL
		result.SkillTags = sql.NullString{String: string(tagsJSON), Valid: true}
	}

	// 3.2 服务区域：非空才更新
	if area := strings.TrimSpace(in.GetServiceArea()); area != "" {
		if len(area) < 2 || len(area) > 64 {
			return nil, errors.New("service_area 长度需在 2~64 之间")
		}
		result.ServiceArea = sql.NullString{String: area, Valid: true}
	}

	// 3.3 状态：0 视为不更新，其它按枚举范围校验
	if s := in.GetStatus(); s != 0 {
		// 3.3.1 校验枚举范围（按你实际业务字典调整，这里假设 1=正常 2=停用 3=审核中）
		if s < 1 || s > 3 {
			return nil, errors.New("status 取值不合法")
		}
		result.Status = int64(s)
	}

	// 3.4 执业证有效期：解析日期 + 校验未过期
	if le := strings.TrimSpace(in.GetLicenseExpire()); le != "" {
		expireTime, err := time.Parse("2006-01-02", le)
		if err != nil {
			return nil, errors.New("license_expire 格式需为 YYYY-MM-DD")
		}
		// 3.4.1 不能设为已过期的日期
		today := time.Now().Truncate(24 * time.Hour)
		if expireTime.Before(today) {
			return nil, errors.New("license_expire 不能早于今天")
		}
		result.LicenseExpire = sql.NullTime{Time: expireTime, Valid: true}
	}

	// ========== 4. 写库 ==========

	// 4.1 全量 Update（result 里其它字段保持原值不变）
	if _, err := l.svcCtx.NurseProfileModel.Update(l.ctx, nil, result); err != nil {
		// 4.2 更新时若有唯一约束冲突（如改手机号等），给友好提示
		// if isDuplicateKeyErr(err) {
		// 	l.Infof("duplicate nurse update, id=%d", in.GetId())
		// 	return nil, errors.New("数据冲突，请检查是否有重复字段")
		// }
		l.Errorf("update nurse failed, id=%d, err=%v", in.GetId(), err)
		return nil, errors.New("更新护士失败")
	}
	ok := true

	// ========== 5. 缓存失效（若 GetNurse 已接缓存，务必加） ==========

	// 5.1 更新成功后清掉详情缓存，避免读到旧数据
	// cacheKey := fmt.Sprintf("nurse:profile:%d", in.GetId())
	// _, _ = l.svcCtx.Redis.Del(cacheKey)

	// ========== 6. 返回 ==========

	return &nurse.UpdateNurseResp{
		Ok: ok,
	}, nil
}
