package domain

import (
	"time"

	// "github.com/kade-chen/library/tools/hash"

	"github.com/kade-chen/google-billing-console/apps/notify"
)

const (
	DEFAULT_DOMAIN = "default"
)

func NewDomain(req *CreateDomainRequest) (*Domain, error) {
	if err := req.validate_strust(); err != nil {
		return nil, err
	}

	// 创建领域对象
	d := &Domain{
		Meta: NewMeta(),
		Spec: req,
	}
	d.Id = req.Name
	// d.Meta.Id = hash.FnvHash(req.Name)
	return d, nil
}

func NewMeta() *Meta {

	return &Meta{
		// Id:       xid.New().String(),
		CreateAt: time.Now().Unix(),
	}
}

// NewCreateDomainRequest todo
func NewCreateDomainRequest() *CreateDomainRequest {
	return &CreateDomainRequest{
		PasswordConfig: NewDefaulPasswordSecurity(),
		LoginSecurity:  NewDefaultLoginSecurity(),
		CodeConfig:     NewDefaultCodeSetting(),
		NotifyConfig:   notify.NewNotifySetting(),
		FeishuSetting:  NewDefaultFeishuConfig(),
		LdapSetting:    NewDefaultLDAPConfig(),
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

// create new default domain
func NewDefaultDomain() *Domain {
	return &Domain{
		Spec: NewCreateDomainRequest(),
	}
}

// NewDescribeDomainRequest 查询详情请求
func NewDescribeDomainRequestByName(name string) *DescribeDomainRequest {
	return &DescribeDomainRequest{
		DescribeBy: DESCRIBE_BY_NAME,
		Name:       name,
	}
}
