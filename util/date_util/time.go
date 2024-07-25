package date_util

import "time"

func TimeNowStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func PreDay() time.Time {
	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	return t.AddDate(0, 0, -1)
}

func PreDayDate() string {
	return PreDay().Format("2006-01-02")
}
func TodayDate() string {
	return time.Now().Format("2006-01-02")
}

func TimeNowNumberStr() string {
	return time.Now().Format("20060102150405")
}
func DateNowHourNumberStr() string {
	return time.Now().Format("2006010215")
}
func DateNowHourMinuteNumberStr() string {
	return time.Now().Format("200601021504")
}
