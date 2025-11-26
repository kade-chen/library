package validate

import (
	"regexp"
	"time"

	"github.com/kade-chen/library/exception"
)

// ValidateYYYYMM 检查日期格式是否为 YYYYMM
func ValidateYYYYMM(value string) error {
	// 正则检查：YYYYMM 且月份 01–12
	re := regexp.MustCompile(`^\d{4}(0[1-9]|1[0-2])$`)
	if !re.MatchString(value) {
		return exception.NewBadRequest("日期格式必须为 YYYYMM，例如 202510")
	}

	// 用 time.Parse 再校验一次
	if _, err := time.Parse("200601", value); err != nil {
		return exception.NewBadRequest("无效的年月")
	}
	return nil
}
