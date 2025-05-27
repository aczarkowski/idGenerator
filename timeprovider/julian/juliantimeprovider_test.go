package julian

import (
	"fmt"
	"strconv"
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

func TestNew(t *testing.T) {
	offset := int64(2000100000)
	provider := New(offset)

	if provider == nil {
		t.Error("Expected provider to be created, got nil")
	}

	if provider.offset != offset {
		t.Errorf("Expected offset %d, got %d", offset, provider.offset)
	}
}

func TestGetTimeStamp(t *testing.T) {
	offset := int64(2000100000)
	provider := New(offset)

	timestamp1 := provider.GetTimeStamp()
	timestamp2 := provider.GetTimeStamp()

	// Timestamps should be close to each other
	diff := timestamp2 - timestamp1
	if diff < 0 || diff > 100 { // Allow small differences due to time progression
		t.Errorf("Timestamps too far apart: %d and %d (diff: %d)", timestamp1, timestamp2, diff)
	}
}

func TestConvertTimeToJulianCalendarLeapYear(t *testing.T) {
	// Test leap year (2020)
	testedTime := time.Date(2020, 2, 29, 12, 30, 45, 0, time.UTC) // Feb 29th (leap day)
	offset := int64(2000100000)

	// Feb 29th should be day 60 of the year
	expectedDayOfYear := 60
	if testedTime.YearDay() != expectedDayOfYear {
		t.Errorf("Expected day of year %d, got %d", expectedDayOfYear, testedTime.YearDay())
	}

	actual := convertTimeToJulianCalendarIncludingTime(testedTime, offset)

	// Should be positive after subtracting offset
	if actual <= 0 {
		t.Errorf("Expected positive result after offset, got %d", actual)
	}
}

func TestConvertTimeToJulianCalendarEndOfYear(t *testing.T) {
	// Test December 31st
	testedTime := time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC)
	offset := int64(2000100000)

	// Dec 31st should be day 365 of the year (2021 is not a leap year)
	expectedDayOfYear := 365
	if testedTime.YearDay() != expectedDayOfYear {
		t.Errorf("Expected day of year %d, got %d", expectedDayOfYear, testedTime.YearDay())
	}

	actual := convertTimeToJulianCalendarIncludingTime(testedTime, offset)

	// Should be positive after subtracting offset
	if actual <= 0 {
		t.Errorf("Expected positive result after offset, got %d", actual)
	}
}

func TestConvertTimeToJulianCalendarMidnight(t *testing.T) {
	// Test exactly midnight
	testedTime := time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC)
	offset := int64(2000100000)

	actual := convertTimeToJulianCalendarIncludingTime(testedTime, offset)

	// At midnight, seconds since beginning of day should be 0
	// So the last 5 digits should be 00000
	actualWithOffset := actual + offset
	if actualWithOffset%100000 != 0 {
		t.Errorf("Expected last 5 digits to be 00000 for midnight, but got %d", actualWithOffset%100000)
	}
}

func TestConvertTimeToJulianCalendarYearTransition(t *testing.T) {
	// Test year 99 to 00 transition (e.g., 1999 to 2000)
	time1999 := time.Date(1999, 12, 31, 23, 59, 59, 0, time.UTC)
	time2000 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	offset := int64(9900100000) // Adjusted offset for these years

	result1999 := convertTimeToJulianCalendarIncludingTime(time1999, offset)
	result2000 := convertTimeToJulianCalendarIncludingTime(time2000, offset)

	// 1999 should give 99365... and 2000 should give 00001...
	// After subtracting offset, 2000 result should be much smaller
	if result2000 >= result1999 {
		t.Errorf("Expected year 2000 result (%d) to be less than 1999 result (%d) due to year rollover", result2000, result1999)
	}
}

func TestConvertTimeToJulianCalendarEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		time     time.Time
		offset   int64
		expected int64
	}{
		{
			name:     "Year 2000 start",
			time:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			offset:   int64(0),
			expected: int64(100000), // Year 2000 % 100 = 0, Day 1, 0 seconds: "0000100000" -> 100000
		},
		{
			name:     "Leap year day",
			time:     time.Date(2020, 2, 29, 12, 30, 45, 0, time.UTC),
			offset:   int64(0),
			expected: int64(2006045045),
		},
		{
			name:     "End of year",
			time:     time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
			offset:   int64(0),
			expected: int64(2136586399), // Year 21, Day 365, 86399 seconds (23:59:59): "2136586399"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := convertTimeToJulianCalendarIncludingTime(tc.time, tc.offset)
			if actual != tc.expected {
				t.Errorf("For %s: expected %d, got %d", tc.name, tc.expected, actual)
			}
		})
	}
}

func TestJulianTimeFormatConsistency(t *testing.T) {
	// Test that the Julian time format is consistent
	testTime := time.Date(2023, 6, 15, 14, 30, 25, 0, time.UTC)
	offset := int64(0)

	// Calculate expected Julian time using the same string format as the implementation
	year := testTime.Year() % 100 // 23
	dayOfYear := testTime.YearDay() // Should be 166 for June 15
	secondsInDay := testTime.Hour()*3600 + testTime.Minute()*60 + testTime.Second() // 52225

	// Use string format to match the implementation: YY(2) + DDD(3) + SSSSS(5) = 10 digits
	expectedStr := fmt.Sprintf("%02d%03d%05d", year, dayOfYear, secondsInDay)
	expected, _ := strconv.Atoi(expectedStr)
	actual := convertTimeToJulianCalendarIncludingTime(testTime, offset)

	if actual != int64(expected) {
		t.Errorf("Julian time format inconsistency. Expected %d, got %d", expected, actual)
		t.Logf("Year: %d, Day of year: %d, Seconds in day: %d", year, dayOfYear, secondsInDay)
	}
}

func TestNew_JulianTimeProvider(t *testing.T) {
	offset := int64(2000100000)
	provider := New(offset)

	if provider == nil {
		t.Error("Expected provider to be created, got nil")
	}

	if provider.offset != offset {
		t.Errorf("Expected offset %d, got %d", offset, provider.offset)
	}
}

func TestGetTimeStamp_JulianTimeProvider(t *testing.T) {
	offset := int64(2000100000)
	provider := New(offset)

	timestamp := provider.GetTimeStamp()

	if timestamp <= 0 {
		t.Errorf("Expected positive timestamp, got %d", timestamp)
	}
}

func TestConvertTimeToJulianCalendar_EdgeCases(t *testing.T) {
	offset := int64(2000100000)

	// Test year boundary (Dec 31 to Jan 1)
	testCases := []struct {
		name     string
		time     time.Time
		expected int64
	}{
		{
			name:     "Year 2020, Day 1, Start of day",
			time:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: int64(2000100000) - offset,
		},
		{
			name:     "Year 2020, Day 366 (leap year), End of day",
			time:     time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: int64(2036686399) - offset,
		},
		{
			name:     "Year 2021, Day 1, Start of day",
			time:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: int64(2100100000) - offset,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := convertTimeToJulianCalendarIncludingTime(tc.time, offset)
			if actual != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, actual)
			}
		})
	}
}

func TestConvertTimeToJulianCalendar_LeapYear(t *testing.T) {
	offset := int64(2000100000)

	// Test leap year (2020) - Feb 29
	leapYearTime := time.Date(2020, 2, 29, 12, 0, 0, 0, time.UTC)
	result := convertTimeToJulianCalendarIncludingTime(leapYearTime, offset)

	// Day 60 (Feb 29) at 12:00:00 (43200 seconds)
	expected := int64(2006043200) - offset
	if result != expected {
		t.Errorf("Expected %d for leap year Feb 29, got %d", expected, result)
	}
}

func TestConvertTimeToJulianCalendar_TimeComponents(t *testing.T) {
	offset := int64(2000100000)

	testCases := []struct {
		name         string
		hours        int
		minutes      int
		seconds      int
		expectedTime int
	}{
		{"Midnight", 0, 0, 0, 0},
		{"One second", 0, 0, 1, 1},
		{"One minute", 0, 1, 0, 60},
		{"One hour", 1, 0, 0, 3600},
		{"12:30:45", 12, 30, 45, 12*3600 + 30*60 + 45},
		{"23:59:59", 23, 59, 59, 23*3600 + 59*60 + 59},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testTime := time.Date(2020, 1, 1, tc.hours, tc.minutes, tc.seconds, 0, time.UTC)
			result := convertTimeToJulianCalendarIncludingTime(testTime, offset)

			// Extract the time component (last 5 digits)
			timeComponent := result % 100000
			if int(timeComponent) != tc.expectedTime {
				t.Errorf("Expected time component %d, got %d", tc.expectedTime, timeComponent)
			}
		})
	}
}
