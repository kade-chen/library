package user

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
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
		return exception.NewInternalServerError("user or password not error: %s", err.Error())
	}

	return nil
}

// CheckPasswordExpired 检测password是否已经过期
// remindDays 提前多少天提醒用户修改密码
// expiredDays 多少天后密码过期
// BeforeExpiredRemindDays =10  password_expired_days=90
func (p *Password) CheckPasswordExpired(ctx context.Context, remindDays, expiredDays uint, bq *bigquery.Client, userID string) error {
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
		err := p.SetNeedReset(ctx, bq, userID, "Password will expire in %f days, Please reset password", -ex)
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
func (p *Password) SetNeedReset(ctx context.Context, bq *bigquery.Client, userID, format string, a ...interface{}) error {
	sql := fmt.Sprintf(`UPDATE %s
						SET password = (
							SELECT AS STRUCT *
							REPLACE (
								TRUE AS need_reset,
								@reason AS reset_reason,
								@update_at AS update_at
							)
							FROM UNNEST([password])
							)
						WHERE id = @id;
    `, fmt.Sprintf("`%s.%s.%s`", ioc.Config().Get(configs.AppName).(*impl.Service).Default_Project_ID, ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDataset, ioc.Config().Get(configs.AppName).(*impl.Service).GoogleBillingConsoleDatasetTableUser))

	q := bq.Query(sql)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: userID},
		{Name: "reason", Value: fmt.Sprintf(format, a...)},
		{Name: "update_at", Value: time.Now().Unix()},
	}

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if status.Err() != nil {
		return status.Err()
	}
	return exception.NewPasswordExired("Update user password configuration failed, %v", nil)
}
