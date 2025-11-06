package tools

import (
	"time"
)

func PartitionTime(StartDate, EndDate string) (PartitionStartTime, PartitionEndTime string) {
	// 解析成 time.Time
	layout := "2006-01-02"
	start, _ := time.Parse(layout, StartDate)
	end, _ := time.Parse(layout, EndDate)

	return start.AddDate(0, 0, -3).Format(layout), end.AddDate(0, 0, 3).Format(layout)
}
