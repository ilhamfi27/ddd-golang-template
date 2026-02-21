package utils

import (
	"fmt"
	"time"
)

// parseTime parses a string into a *time.Time
func ParseTime(dateStr string) (*time.Time, error) {
	layout := "2006-01-02" // adjust the layout to match your date format
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetNow() time.Time {
	return time.Now()
}

func CalculateDuration(fromDate string) (int, int, int) {
	// Parse the input date
	startDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		panic(fmt.Sprintf("Invalid date format: %v", err))
	}

	// Get the current date
	currentDate := GetNow()

	// Calculate the total days between the two dates
	duration := currentDate.Sub(startDate)
	totalDays := int(duration.Hours() / 24)

	// Define basic constants
	daysInYear := 365
	daysInMonth := 30

	// Calculate years
	years := totalDays / daysInYear
	remainingDays := totalDays % daysInYear

	// Calculate months
	months := remainingDays / daysInMonth
	days := remainingDays % daysInMonth

	return years, months, days
}
