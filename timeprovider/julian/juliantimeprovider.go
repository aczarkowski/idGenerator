package julian

import (
	"fmt"
	"strconv"
	"time"
)

// TimeProvider Implement the TimeProvider interface using Julian calendar
type timeProvider struct {
	offset int64
}

// New TimeProvider
func New(offset int64) *timeProvider {
	return &timeProvider{
		offset: offset,
	}
}

// GetTimeStamp Returns the current time in milliseconds
func (t *timeProvider) GetTimeStamp() int64 {
	return convertTimeToJulianCalendarIncludingTime(time.Now().UTC(), t.offset)
}

func convertTimeToJulianCalendarIncludingTime(time time.Time, offset int64) int64 {
	// From time get last 2 digits of the year
	last2DigitsOfTheYear := time.Year() % 100

	// From time calculate number of days since beginning of the year
	daysSinceBeginningOfTheYear := time.YearDay()

	// From time get number of seconds since beginning of the current day
	secondsSinceBeginningOfTheDay := time.Hour()*3600 + time.Minute()*60 + time.Second()

	// Concatenate all the above into a string and convert to a number

	julianTime, _ := strconv.Atoi(fmt.Sprintf("%d%03d%05d", last2DigitsOfTheYear, daysSinceBeginningOfTheYear, secondsSinceBeginningOfTheDay))

	return int64(julianTime) - offset
}
