/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/25 14:28
 */
package time_utils

import (
	"testing"
	"time"
)

// TestSetUnixTime tests the SetUnixTime function.
func TestSetUnixTime(t *testing.T) {
	tests := []struct {
		name     string
		source   time.Time
		expected int64
	}{
		{
			name:     "non-zero time",
			source:   time.Date(2023, 10, 5, 14, 30, 0, 0, time.UTC),
			expected: 1696516200,
		},
		{
			name:     "zero time",
			source:   time.Time{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target int64
			SetUnixTime(&target, tt.source)

			if target != tt.expected {
				t.Errorf("SetUnixTime() = %v, want %v", target, tt.expected)
			}
		})
	}
}

func TestFormatYMW_1(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected [3]int32
	}{
		{
			name:     "2024-04-05 Normal Date",
			input:    time.Date(2024, time.April, 5, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2024, 4, 14},
		},
		{
			name:     "2023-12-31 Year Boundary (Next Year Week)",
			input:    time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2023, 12, 52},
		},
		{
			name:     "2024-01-01 New Year First Day",
			input:    time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2024, 1, 1},
		},
		{
			name:     "2023-01-01 New Year First Day (Non-Leap)",
			input:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2022, 1, 52},
		},
		{
			name:     "2020-12-31 Last Day of 53 Week Year",
			input:    time.Date(2020, time.December, 31, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2020, 12, 53},
		},
		{
			name:     "2024-01-07 Cross-Week Transition (Belongs to 2023W52)",
			input:    time.Date(2024, time.January, 7, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2024, 1, 1},
		},
		{
			name:     "2023-02-28 Last Day of February (Non-Leap)",
			input:    time.Date(2023, time.February, 28, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2023, 2, 9},
		},
		{
			name:     "2024-02-29 Last Day of February (Leap Year)",
			input:    time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2024, 2, 9},
		},
		{
			name:     "2023-03-01 First Day of March",
			input:    time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2023, 3, 9},
		},
		{
			name:     "2024-12-31 Last Day of Year (May belong to next year week)",
			input:    time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
			expected: [3]int32{2025, 12, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, week := FormatYearMonthWeek(tt.input)
			if year != tt.expected[0] || month != tt.expected[1] || week != tt.expected[2] {
				t.Errorf("FormatYMW(%v) = (%d, %d, %d), want (%d, %d, %d)",
					tt.input, year, month, week,
					tt.expected[0], tt.expected[1], tt.expected[2])
			}
		})
	}
}
