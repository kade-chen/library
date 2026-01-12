package organization

import (
	"time"

	// "github.com/kade-chen/library/tools/hash"

	"github.com/kade-chen/google-billing-console/apps/notify"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/tools/hash"
)

const (
	DEFAULT_DOMAIN = "default"
)

func NewOrganization(req *CreateOrganizationRequest) (*Organization, error) {
	if err := req.validate_strust(); err != nil {
		return nil, exception.NewBadRequest("invalid request, ERROR: %s", err)
	}

	// 创建领域对象
	d := &Organization{
		Meta: NewMeta(),
		Spec: req,
	}
	d.Id = req.OrganizationDetail.SubOrganization
	d.Id = hash.FnvHash(req.OrganizationDetail.SubOrganization)
	return d, nil
}

func NewMeta() *Meta {

	return &Meta{
		// Id:       xid.New().String(),
		CreateAt: time.Now().Unix(),
	}
}

// NewCreateOrganizationRequest todo
func NewCreateOrganizationRequest() *CreateOrganizationRequest {
	return &CreateOrganizationRequest{
		PasswordConfig:  NewDefaulPasswordSecurity(),
		LoginSecurity:   NewDefaultLoginSecurity(),
		CodeConfig:      NewDefaultCodeSetting(),
		NotifyConfig:    notify.NewNotifySetting(),
		FeishuSetting:   NewDefaultFeishuConfig(),
		LdapSetting:     NewDefaultLDAPConfig(),
		DingdingSetting: &DingDingConfig{},
		Contack:         &Contact{},
		WechatWorkSetting: &WechatWorkConfig{
			AccessToken: &WechatWorkAccessToken{},
		},
		OrganizationDetail: &OrganizationDetail{},
	}
}

// NewDefaulPasswordSecurity todo
func NewDefaulPasswordSecurity() *PasswordConfig {
	return &PasswordConfig{
		Enabled:                 true,
		Length:                  8,
		IncludeNumber:           true,
		IncludeLowerLetter:      true,
		IncludeUpperLetter:      false,
		IncludeSymbols:          false,
		RepeateLimite:           1,
		PasswordExpiredDays:     90,
		BeforeExpiredRemindDays: 10,
	}
}

// NewDefaultLoginSecurity todo
func NewDefaultLoginSecurity() *LoginSecurity {
	return &LoginSecurity{
		ExceptionLock: false,
		ExceptionLockConfig: &ExceptionLockConfig{
			OtherPlaceLogin: true,
			NotLoginDays:    30,
		},
		RetryLock: true,
		RetryLockConfig: &RetryLockConfig{
			RetryLimite:  5,
			LockedMinite: 30,
		},
		IpLimite: false,
		IpLimiteConfig: &IPLimiteConfig{
			Ip: []string{},
		},
	}
}

// NewDefaultCodeSetting todo
func NewDefaultCodeSetting() *CodeSetting {
	return &CodeSetting{
		NotifyType:    notify.NOTIFY_TYPE_MAIL,
		ExpireMinutes: 10,
		MailTemplate:  "您的动态验证码为：{1}，{2}分钟内有效！，如非本人操作，请忽略本邮件！",
	}
}

// create new default Organization
func NewDefaultOrganization() *Organization {
	return &Organization{
		Spec: NewCreateOrganizationRequest(),
	}
}

// NewDescribeOrganizationRequest 查询详情请求
func NewDescribeOrganizationRequestByName(name string) *DescribeOrganizationRequest {
	return &DescribeOrganizationRequest{
		DescribeBy: DESCRIBE_BY_NAME,
		Name:       name,
	}
}
