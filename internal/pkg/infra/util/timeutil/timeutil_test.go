package timeutil

import (
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {
	result := GetCurrentFormattedTime()
	tm, err := time.Parse(TIME_FORMAT, result)

	if err != nil {
		t.Errorf("Parsing failed: %v", err)
	}

	if tm.Format(TIME_FORMAT) != result {
		t.Errorf("Unexpected error in GetCurrentFormattedTime, got %v, but expected %v.", tm.Format(TIME_FORMAT), result)
	}
}

func TestConvertEpochToFormattedTime(t *testing.T) {
	// 2020-10-10 12:00:00 Shanghai
	epochMS := int64(1602302400000)

	result := ConvertEpochMSToFormattedTime(epochMS)

	expected := "2020-10-10 12:00:00"
	if result != expected {
		t.Errorf("Unexpected error in ConvertEpochMSToFormattedTime, got %v, but expected %v.", result, expected)
	}
}
