package user

import (
	"time"

	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/tools/format"
)

func NewUser(req *CreateUserRequest) (*User, error) {
	if req == nil {
		u := &User{
			Meta: &Meta{},
			Spec: &CreateUserRequest{
				Labels:     map[string]string{},
				Feishu:     &Feishu{},
				Dingding:   &DingDing{},
				Wechatwork: &WechatWork{},
				Profile:    &Profile{},
			},
			Password:      &Password{},
			Status:        &Status{},
			FeishuToken:   &FeishuAccessToken{},
			DingdingToken: &DingDingAccessToken{},
		}
		return u, nil
	} else {
		// 1.generate password strust
		password, err := NewHashPassword(req.Password)
		if err != nil {
			return nil, exception.NewBadRequest(err.Error())
		}
		// 2.generate user strust
		u := &User{
			Meta: &Meta{
				// Id:       req.Username + "-" + xid.New().String(),
				Id:       req.Username + "@" + req.Domain,
				CreateAt: time.Now().Unix(),
			},
			Spec:     req,
			Password: password,
			Status: &Status{
				IsInitialized: false,
				Locked:        false,
			},
			FeishuToken:   &FeishuAccessToken{},
			DingdingToken: &DingDingAccessToken{},
		}
		// 3.validate strust
		if err := validate.Struct(u); err != nil {
			return nil, exception.NewBadRequest("validate error: %s", err.Error())
		}
		return u, nil
	}
}

func NewDefaultUser() *User {
	return &User{}
}

// initialize some null value,compatible with provious data
func (u *User) InitializeNullValue() {
	if u.FeishuToken == nil {
		u.FeishuToken = &FeishuAccessToken{
			IssueAt: time.Now().Unix(),
		}
	}
	if u.Spec.Feishu == nil {
		u.Spec.Feishu = &Feishu{}
	}
	if u.Spec.Dingding == nil {
		u.Spec.Dingding = &DingDing{}
	}
	if u.Spec.Wechatwork == nil {
		u.Spec.Wechatwork = &WechatWork{}
	}
}

// desensitization user password and history
func (u *User) Desensitization() {
	if u.Password != nil {
		u.Password.Password = "******"
		u.Password.History = []string{"******"}
	}
}

func (u *User) ToJson() string {
	return format.ToJSON(u)
}
