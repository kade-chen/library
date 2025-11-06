package format

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

func FormatAny(n any) any {
	switch v := n.(type) {
	case bigquery.NullFloat64:
		if v.Valid {
			return fmt.Sprintf("%.2f", v.Float64) // 保留两位小数
		}
	case bigquery.NullString:
		if v.Valid {
			return v.String() // 注意大括号：调用方法
		}
	case bigquery.NullDate:
		if v.Valid {
			return v.Date // 转成 yyyy-mm-dd
		}
	case float64:
		return fmt.Sprintf("%.2f", v)
	case int64:
		return fmt.Sprintf("%2d", v)
	case string:
		return v
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
	return n
}
