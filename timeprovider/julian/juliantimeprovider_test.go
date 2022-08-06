package julian

import (
	"testing"
	"time"
)

func TestConvertTimeToJulianCalendarIncludingTime(t *testing.T) {
	testedTime := time.Date(2020, 10, 12, 13, 14, 15, 0, time.UTC)
	offset := int64(2000100000)
	dayNumber := testedTime.YearDay()
	t.Logf("Day number: %d", dayNumber)
	expected := int64(2028647655) - offset
	actual := convertTimeToJulianCalendarIncludingTime(testedTime, offset)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestConvertTimeToJulianCalendarFor20220101(t *testing.T) {
	testedTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	offset := int64(2000100000)
	expected := int64(2200100000) - offset
	actual := convertTimeToJulianCalendarIncludingTime(testedTime, offset)
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
