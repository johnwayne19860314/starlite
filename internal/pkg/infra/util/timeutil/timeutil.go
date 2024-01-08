package timeutil

import (
	"fmt"
	"time"

	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

// GetCurrentFormattedTime yyyy-MM-dd HH:mm:ss
func GetCurrentFormattedTime() string {
	t := time.Now().In(getDefaultTimeLocation())
	return t.Format(TIME_FORMAT)
}

// GetCurrentEpochMSTime 1691136749245
func GetCurrentEpochMSTime() int64 {
	return time.Now().In(getDefaultTimeLocation()).UnixMilli()
}

// ConvertEpochMSToFormattedTime convert Epoch ms time into formatted time yyyy-MM-dd HH:mm:ss
func ConvertEpochMSToFormattedTime(epochMS int64) string {
	// convert to formatted time
	t := time.UnixMilli(epochMS).In(getDefaultTimeLocation())
	return t.Format(TIME_FORMAT)
}

// ConvertRFC3339ToFormattedTime convert 2023-07-03T12:18:08Z into yyyy-MM-dd HH:mm:ss
func ConvertRFC3339ToFormattedTime(rfcTime string) string {
	parsedTime, err := time.Parse(time.RFC3339, rfcTime)
	if err != nil {
		return rfcTime
	}
	return parsedTime.Format(TIME_FORMAT)
}

func getDefaultTimeLocation() *time.Location {
	loc, err := time.LoadLocation("Shanghai")
	if err != nil {
		logx.Error(fmt.Sprintf("failed to get timeLocation for %v", "shanghai"))
		return time.UTC
	}
	return loc
}
