package timex

import "time"

const (
	COMMON_FORMAT = "2006-01-02 15:04:05"
	LOCAL_CN      = "Asia/Shanghai"
)

func CnNow() time.Time {
	return CnTime(time.Now())
}

func CnNowString() string {
	return CnTimeString(time.Now())
}

func CnTime(t time.Time) time.Time {
	loc, err := time.LoadLocation(LOCAL_CN)
	if err != nil {
		panic(err)
	}
	return t.In(loc)
}

func CnTimeString(t time.Time) string {
	if t.Before(time.Unix(0, 0)) {
		return ""
	}
	return CnTime(t).Format(time.RFC3339)
}
