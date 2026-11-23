package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"Nurse/app/nurse/cmd/rpc/internal/svc"
	"Nurse/app/nurse/cmd/rpc/nurse"
	"Nurse/app/nurse/model"

	"github.com/zeromicro/go-zero/core/logx"
)

// 预编译手机号正则：中国大陆 11 位手机号
var mobilePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

// 预编译执业证编号正则：示例规则，可按业务替换
// 这里假设：字母+数字组合，长度 6~32
var licenseNoPattern = regexp.MustCompile(`^[A-Za-z0-9]{6,32}$`)

type CreateNurseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNurseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNurseLogic {
	return &CreateNurseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNurseLogic) CreateNurse(in *nurse.CreateNurseReq) (*nurse.CreateNurseResp, error) {
	// todo: add your logic here and delete this line

	// ========== 1. 参数基础校验 ==========

	// 1.1 校验真实姓名：非空，去掉首尾空白后长度 2~32
	realName := strings.TrimSpace(in.GetRealName())
	if len(realName) < 2 || len(realName) > 32 {
		return nil, errors.New("real_name 长度需在 2~32 之间")
	}

	// 1.2 校验手机号：必须匹配 11 位大陆手机号
	mobile := strings.TrimSpace(in.GetMobile())
	if !mobilePattern.MatchString(mobile) {
		return nil, errors.New("mobile 格式不合法")
	}

	// 1.3 校验执业证编号：非空且符合格式
	licenseNo := strings.TrimSpace(in.GetLicenseNo())
	if !licenseNoPattern.MatchString(licenseNo) {
		return nil, errors.New("license_no 格式不合法")
	}

	// 1.4 校验执业证有效期：必须为 YYYY-MM-DD 格式，且不能早于今天
	licenseExpire := strings.TrimSpace(in.GetLicenseExpire())
	expireTime, err := time.Parse("2006-01-02", licenseExpire)
	if err != nil {
		return nil, errors.New("license_expire 格式需为 YYYY-MM-DD")
	}
	// 用当天 0 点做比较，避免时分秒影响
	today := time.Now().Truncate(24 * time.Hour)
	if expireTime.Before(today) {
		return nil, errors.New("license_expire 已过期")
	}

	// 1.5 校验服务区域：非空，长度 2~64
	serviceArea := strings.TrimSpace(in.GetServiceArea())
	if len(serviceArea) < 2 || len(serviceArea) > 64 {
		return nil, errors.New("service_area 长度需在 2~64 之间")
	}

	// 1.6 校验技能标签：数量限制 + 单个标签去空
	if len(in.GetSkillTags()) > 20 {
		return nil, errors.New("skill_tags 数量不能超过 20 个")
	}
	skillTags := make([]string, 0, len(in.GetSkillTags()))
	for _, tag := range in.GetSkillTags() {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue // 跳过空标签
		}
		skillTags = append(skillTags, tag)
	}
	// 若传入但全为空，则等同于未传；可按需改成报错

	// ========== 2. 业务唯一性校验 ==========

	// 2.1 手机号必须全局唯一（一个手机号只能注册一个护士）
	exists, err := l.svcCtx.NurseProfileModel.FindOneByMobile(l.ctx, mobile)
	if err != nil {
		l.Errorf("check mobile exists failed, mobile=%s, err=%v", mobile, err)
		return nil, errors.New("系统繁忙，请稍后重试")
	}
	if exists.Mobile == mobile {
		return nil, errors.New("该手机号已被注册")
	}

	// 2.2 执业证编号也必须全局唯一
	exists, err = l.svcCtx.NurseProfileModel.FindOneByLicenseNo(l.ctx, licenseNo)
	if err != nil {
		l.Errorf("check license_no exists failed, license_no=%s, err=%v", licenseNo, err)
		return nil, errors.New("系统繁忙，请稍后重试")
	}
	if exists.LicenseNo == licenseNo {
		return nil, errors.New("该执业证编号已被注册")
	}

	// ========== 3. 数据组装 ==========

	// 3.1 技能标签序列化为 JSON 字符串，方便 DB 存储与后续检索
	skillTagsJSON, err := json.Marshal(skillTags)
	if err != nil {
		l.Errorf("marshal skill_tags failed, err=%v", err)
		return nil, errors.New("系统繁忙，请稍后重试")
	}

	// 3.2 组装要插入 DB 的实体
	nurseData := &model.NurseProfile{
		RealName:  realName,
		Mobile:    mobile,
		LicenseNo: licenseNo,
		LicenseExpire: sql.NullTime{ // ← 修改点：由 time.Time 改为 sql.NullTime
			Time:  expireTime,
			Valid: true,
		}, // DB 用 DATE / DATETIME 类型
		SkillTags: sql.NullString{
			String: string(skillTagsJSON),
			Valid:  true, // 表示非 NULL
		}, // DB 用 VARCHAR / TEXT 类型
		ServiceArea: sql.NullString{ // ← 修改点：string → sql.NullString
			String: serviceArea,
			Valid:  true, // 表示非 NULL
		},
		Status: 1, // 假设 1=正常，可按需扩展：待审核/在职等
	}

	// 4.1 插入记录（model 内部处理主键回填）
	result, err := l.svcCtx.NurseProfileModel.Insert(l.ctx, nil, nurseData)
	if err != nil {
		// 4.2 针对唯一索引冲突给出友好提示
		// （防止并发下第 2 步的校验被绕过）
		// if model.IsDuplicateKeyErr(err) {
		// 	l.Infof("duplicate nurse insert, mobile=%s, license_no=%s", mobile, licenseNo)
		// 	return nil, errors.New("手机号或执业证编号已存在")
		// }
		// l.Errorf("insert nurse failed, data=%+v, err=%v", nurseData, err)
		return nil, errors.New("创建护士失败")
	}

	id, err := result.LastInsertId()
	if err != nil {

		return nil, errors.New("创建护士失败")
	}

	return &nurse.CreateNurseResp{
		Id: id,
	}, nil
}
