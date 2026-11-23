package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetNurseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNurseLogic {
	return &GetNurseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNurseLogic) GetNurse(in *nurse.GetNurseReq) (*nurse.GetNurseResp, error) {
	// todo: add your logic here and delete this line

	// 1.1 ID 必须是正整数，否则直接拒绝，避免无效 DB 查询
	if in.GetId() <= 0 {
		return nil, errors.New("id 不合法")
	}

	result, err := l.svcCtx.NurseProfileModel.FindOne(l.ctx, uint64(in.Id))

	if err != nil {
		// 2.2 区分"记录不存在"和"其它 DB 错误"
		//     sqlx.ErrNotFound 是 go-zero sqlx 包定义的哨兵错误
		if errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("护士不存在")
		}
		// 2.3 其它错误打日志并返回通用提示，避免泄露内部信息
		l.Errorf("find nurse profile failed, id=%d, err=%v", in.GetId(), err)
		return nil, errors.New("系统繁忙，请稍后重试")
	}

	// 2.4 防御性判空（理论上 FindOne 成功时不会是 nil，但避免 panic）
	if result == nil {
		return nil, errors.New("护士不存在")
	}

	var skillTags []string
	if result.SkillTags.Valid && result.SkillTags.String != "" {
		if err := json.Unmarshal([]byte(result.SkillTags.String), &skillTags); err != nil {
			// 3.2 反序列化失败不阻断主流程，记日志后降级为空切片
			l.Errorf("unmarshal skill_tags failed, id=%d, raw=%s, err=%v",
				in.GetId(), result.SkillTags.String, err)
			skillTags = []string{}
		}
	}
	if skillTags == nil {
		skillTags = []string{} // 保证接口返回 [] 而非 null，前端更好处理
	}

	// 3.3 执业证有效期：sql.NullTime → "2006-01-02" 字符串（Null 时留空）
	var licenseExpire string
	if result.LicenseExpire.Valid {
		licenseExpire = result.LicenseExpire.Time.Format("2006-01-02")
	}

	var serviceArea string
	if result.ServiceArea.Valid {
		serviceArea = result.ServiceArea.String
	}

	info := &nurse.NurseInfo{
		Id:            int64(result.Id),
		RealName:      result.RealName,
		Mobile:        desensitizeMobile(result.Mobile),
		LicenseNo:     result.LicenseNo,
		LicenseExpire: licenseExpire,
		SkillTags:     skillTags,
		ServiceArea:   serviceArea,
		Status:        int32(result.Status),
	}

	return &nurse.GetNurseResp{
		Info: info,
	}, nil
}

// desensitizeMobile 手机号脱敏：138****8000
func desensitizeMobile(mobile string) string {
	if len(mobile) != 11 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}
