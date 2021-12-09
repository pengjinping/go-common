package timehelper

import (
	"fmt"
	"time"
)

func FormatDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}
func FormatTime(t time.Time) string {
	hour, min, sec := t.Clock()
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}
func FormatDateTime(t time.Time) string {
	return fmt.Sprintf("%s %s", FormatDate(t), FormatTime(t))
}
