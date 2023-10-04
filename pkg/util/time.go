package util

import (
	"fmt"
	"regexp"
	"time"
)

func ResolveTime(startTimeStamp, endTimeStamp int64) string {
	// 将 Unix 时间戳转换为 time.Time 类型
	startTime := time.Unix(startTimeStamp, 0)
	endTime := time.Unix(endTimeStamp, 0)

	// 计算时间差
	duration := endTime.Sub(startTime)

	// 格式化时间差为字符串
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	timeStr := fmt.Sprintf("Use: %d Hour %d Minute %d Second", hours, minutes, seconds)

	return timeStr
}

// GetTimestampAfterDay 获取几天后的时间戳
func GetTimestampAfterDay(days int) int64 {
	now := time.Now()
	duration := time.Duration(days) * 24 * time.Hour
	futureTime := now.Add(duration)
	return futureTime.Unix()
}

// FormatTimestampToDateTime 时间戳格式化为日期
func FormatTimestampToDateTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// GetDayTimestampZ 取当天凌晨0点时间戳
func GetDayTimestampZ() int64 {
	// 获取当前时间
	currentTime := time.Now()
	// 获取年、月、日
	year, month, day := currentTime.Date()
	// 创建当天凌晨0点的时间对象
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, currentTime.Location())
	// 获取时间戳
	return startOfDay.Unix()
}

// GetLogTimestampZ 获取 "Sep  16 15:44:54 "或"Sep  6 15:44:54 "字符串格式的时间戳
func GetLogTimestampZ(log string) int64 {
	re := regexp.MustCompile(`(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})`)
	match := re.FindStringSubmatch(log)
	if len(match) > 1 {
		timestamp, err := time.Parse("2006 Jan  2 15:04:05", fmt.Sprintf("%d %s", time.Now().Year(), match[1]))
		if err == nil {
			return timestamp.Unix()
		}
	}
	return 0
}

// HoursBetweenTimestamps 返回两个时间戳内的年月日(string) - 小时([]int)
func HoursBetweenTimestamps(startTimestamp, endTimestamp int64) (map[string][]int, error) {
	if startTimestamp > endTimestamp {
		return nil, fmt.Errorf("startTimestamp should be less than endTimestamp")
	}

	hoursMap := make(map[string][]int)
	startTime := time.Unix(startTimestamp, 0)
	endTime := time.Unix(endTimestamp, 0)

	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(time.Hour) {
		dateKey := currentTime.Format("20060102")
		hour := currentTime.Hour()

		if _, exists := hoursMap[dateKey]; !exists {
			hoursMap[dateKey] = []int{hour}
		} else {
			hoursMap[dateKey] = append(hoursMap[dateKey], hour)
		}
	}

	return hoursMap, nil
}
