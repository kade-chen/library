package time

import (
	"fmt"
	"time"
)

func PartitionTime(StartDate, EndDate string) (PartitionStartTime, PartitionEndTime string) {
	// 解析成 time.Time
	layout := "2006-01-02"
	start, _ := time.Parse(layout, StartDate)
	end, _ := time.Parse(layout, EndDate)
	return start.AddDate(0, 0, -2).Format(layout), end.AddDate(0, 0, 2).Format(layout)
}

func CustomPartitionTime(StartDate, EndDate string, days int) (PartitionStartTime, PartitionEndTime string) {
	// 解析成 time.Time
	layout := "2006-01-02"
	start, _ := time.Parse(layout, StartDate)
	end, _ := time.Parse(layout, EndDate)
	return start.AddDate(0, 0, -days).Format(layout), end.AddDate(0, 0, days).Format(layout)
}

// func InvoiceMonthPartitionTime(startStr, endStr string) (string, string, error) {
// 	// 1. 解析 YYYYMM → time.Time（默认取每月1号）
// 	start, err := time.Parse("200601", startStr)
// 	if err != nil {
// 		return "", "", fmt.Errorf("start_date invalid: %w", err)
// 	}
// 	end, err := time.Parse("200601", endStr)
// 	if err != nil {
// 		return "", "", fmt.Errorf("end_date invalid: %w", err)
// 	}

// 	// 2. 计算 month_span = DATE_DIFF(MONTH) + 1
// 	yearDiff := end.Year() - start.Year()
// 	monthDiff := int(end.Month()) - int(start.Month())
// 	monthSpan := yearDiff*12 + monthDiff + 1

// 	// 3. 前一区间：start - span, end - span
// 	prevStart := start.AddDate(0, -monthSpan, 0)
// 	prevEnd := end.AddDate(0, -monthSpan, 0)

// 	// 4. 返回 YYYY-MM-01（模拟 BigQuery PARSE_DATE(%Y%m) 输出）
// 	layout := "2006-01-02"
// 	return prevStart.Format(layout), prevEnd.Format(layout), nil
// }

// PartitionPrevDates 根据 startDate 和 endDate，计算 prevStart 和 prevEnd
func PartitionPrevDates(startDateStr, endDateStr string) (prevStart, prevEnd string) {
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, startDateStr)
	// if err != nil {
	// 	return "", "", fmt.Errorf("解析 startDate 失败: %w", err)
	// }
	endDate, _ := time.Parse(layout, endDateStr)
	// if err != nil {
	// 	return "", "", fmt.Errorf("解析 endDate 失败: %w", err)
	// }
	// 计算日期差
	diffDays := int(endDate.Sub(startDate).Hours() / 24)
	// SQL 中 prev_start = DATE_SUB(start_date, INTERVAL DATE_DIFF + 1 DAY)
	// prevStartTime := startDate.AddDate(0, 0, -(diffDays + 1)).Format(layout)
	// prevEndTime := endDate.AddDate(0, 0, -(diffDays + 1)).Format(layout)
	return startDate.AddDate(0, 0, -(diffDays + 1)).Format(layout), endDate.AddDate(0, 0, -(diffDays + 1)).Format(layout)
}

// 202508", "202510"
// 202505 202507
func InvoiceMonthTime(startStr, endStr string) (string, string) {
	// 1. 解析 YYYYMM → time.Time（默认取每月1号）
	start, err := time.Parse("200601", startStr)
	if err != nil {
		return "", ""
	}

	end, err := time.Parse("200601", endStr)
	if err != nil {
		return "", ""
	}

	// 2. 计算 month_span = DATE_DIFF(MONTH) + 1
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())
	monthSpan := yearDiff*12 + monthDiff + 1

	// 3. 前一区间：start - span, end - span
	prevStart := start.AddDate(0, -monthSpan, 0)
	prevEnd := end.AddDate(0, -monthSpan, 0)

	// 4. 返回 YYYYMM
	return prevStart.Format("200601"), prevEnd.Format("200601")
}

func InvoiceMonthPartitionTime(startStr, endStr string) (string, string) {
	// 解析 YYYYMM → time.Time（默认取每月1号）
	start, err := time.Parse("200601", startStr)
	if err != nil {
		return "", ""
	}

	end, err := time.Parse("200601", endStr)
	if err != nil {
		return "", ""
	}

	// prev_start = start - 1 个月（取每月 1 号）
	prevStart := start.AddDate(0, -1, 0)
	prevStartStr := fmt.Sprintf("%04d-%02d-01", prevStart.Year(), prevStart.Month())

	// prev_end = end - 1 个月的“月底”
	// 技巧：跳到下个月 1 号再往前减一天，就是这个月的月底
	prevEndMonth := end.AddDate(0, 1, 0)
	nextMonth := time.Date(prevEndMonth.Year(), prevEndMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	prevEnd := nextMonth.Add(-24 * time.Hour)
	prevEndStr := prevEnd.Format("2006-01-02")

	return prevStartStr, prevEndStr
}
