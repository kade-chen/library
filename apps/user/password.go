package user

import (
	"context"
	"fmt"
	"time"

	"github.com/kade-chen/library/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// hash password
func NewHashPassword(password string) (*Password, error) {
	//1.generate hash password
	bycrptpassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, exception.NewBadRequest("hash password error: %s", err.Error())
	}
	//2.create password object
	p := &Password{
		Password:      string(bycrptpassword),
		CreateAt:      time.Now().Unix(),
		UpdateAt:      time.Now().Unix(),
		ExpiredDays:   180,
		ExpiredRemind: 30,
	}
	return p, nil
}

// CheckPassword 判断password 是否正确
func (p *Password) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	if err != nil {
		return exception.NewBadRequest("用户名或者密码不正确")
	}

	return nil
}

// CheckPasswordExpired 检测password是否已经过期
// remindDays 提前多少天提醒用户修改密码
// expiredDays 多少天后密码过期
// BeforeExpiredRemindDays =10  password_expired_days=90
func (p *Password) CheckPasswordExpired(remindDays, expiredDays uint, col *mongo.Collection, userID string) error {
	// 永不过期
	if expiredDays == 0 {
		return nil
	}

	now := time.Now()
	expiredAt := time.Unix(p.CreateAt, 0).Add(time.Duration(expiredDays) * time.Hour * 24)
	//没过期是负值
	ex := now.Sub(expiredAt).Hours() / 24
	// 提前提醒10天
	if 0 < -ex && -ex <= float64(remindDays) {
		err := p.SetNeedReset(col, userID, "Password will expire in %f days, Please reset password", -ex)
		return err
	}

	// if ex >= -float64(remindDays) {
	// 	p.SetNeedReset(col, userID, "密码%f天后过期, 请重置密码", -ex)
	// }
	if ex > 0 {
		return exception.NewPasswordExired("password expired %f days, Please contact the administrator to reset the password", ex)
	}

	return nil
}

// SetNeedReset 需要被重置
func (p *Password) SetNeedReset(col *mongo.Collection, userID, format string, a ...interface{}) error {

	_, err := col.UpdateOne(context.TODO(), bson.M{"_id": userID}, bson.M{
		"$set": bson.M{
			"password.need_reset":   true,
			"password.reset_reason": fmt.Sprintf(format, a...),
			"password.update_at":    time.Now().Unix(),
		},
	})
	return exception.NewPasswordExired("Update user password configuration failed, %v", err)
}
