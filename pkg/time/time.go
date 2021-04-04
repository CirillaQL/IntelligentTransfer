package MyTimeParse

import (
	"IntelligentTransfer/pkg/logger"
	"time"
)

const timeTemplate = "15:04"

// TimeParse 将时间字符串转化为时间格式
func TimeParse(input string) time.Time {
	stamp, err := time.ParseInLocation(timeTemplate, input, time.Local)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Time Parse failed. Time:%+v  Error:%+v", input, err)
		return time.Time{}
	}
	return stamp
}

// TimeCompareLater 比较两个两个时间中较晚的时间
func TimeCompareLater(timeA, timeB string) string {
	stampTimeA := TimeParse(timeA)
	stampTimeB := TimeParse(timeB)
	if stampTimeA.After(stampTimeB) {
		return timeA
	} else {
		return timeB
	}
}
