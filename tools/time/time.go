package format

import (
	"time"
)

func PartitionTime(StartDate, EndDate string) (PartitionStartTime, PartitionEndTime string) {
	// 解析成 time.Time
	layout := "2006-01-02"
	start, _ := time.Parse(layout, StartDate)
	end, _ := time.Parse(layout, EndDate)
	return start.AddDate(0, 0, -30).Format(layout), end.AddDate(0, 0, 30).Format(layout)
}

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
