/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/25 14:26
 */
package time_utils

import (
	"base/consts"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatTimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateTime)
}

func SetUnixTime(target *int64, source time.Time) {
	if !source.IsZero() {
		*target = source.Unix()
	}
}

func SetStringTime(target *string, source time.Time) {
	if !source.IsZero() {
		*target = source.Format(time.DateTime)
	}
}

func TimeToString(source time.Time) string {
	if !source.IsZero() {
		return source.Format(time.DateTime)
	}
	return ""
}

func FormatYearMonthWeek(source time.Time) (int32, int32, int32) {
	y, w := source.ISOWeek()
	return int32(y), int32(source.Month()), int32(w)
	//year := source.Year()
	//month := int32(source.Month())
	//
	//yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, source.Location())
	//dayOfYear := source.YearDay() - 1 // 从 0 开始
	//
	//firstDayOfWeek := (int(yearStart.Weekday()) + 6) % 7 // 周一为每周第一天
	//week := (dayOfYear + firstDayOfWeek) / 7
	//
	//return int32(year), month, int32(week + 1)
}

func FormatYMDhmss() string {
	timeStr := time.Now().UTC().Format(consts.TIME_LAYOUT_FORMAT_YMDHMSS)
	replacements := []string{"-", ":", ".", " "}
	replacement := ""

	processedString := timeStr
	for _, char := range replacements {
		processedString = strings.ReplaceAll(processedString, char, replacement)
	}
	return processedString
}

func FormatMMDDHHMMSS() string {
	return time.Now().UTC().Format(consts.TIME_LAYOUT_FORMAT_YYMMDDHHMMSS)

}

// ConvertTimeToCurrency retrieves the current time in the timezone
// corresponding to the given currency code.
//
// This method internally obtains the current UTC time and then converts it
// to the timezone associated with the specified currency.
//
// Parameters:
//
//	currencyCode string: The currency code (e.g., "USD", "EUR").
//
// Returns:
//
//	time.Time: The converted current local time in the currency's timezone.
//	           If the currency is not found or the timezone fails to load,
//	           the current UTC time is returned.
//	error: An error if the currency's timezone information is not found
//	       or if the corresponding timezone fails to load.
func ConvertTimeToCurrency(currencyCode string) (time.Time, error) {
	// 1. Obtain the current UTC time within the method.
	//    time.Now() returns the local time, and .UTC() converts it to Coordinated Universal Time.
	currentUTCTime := time.Now().UTC()

	// 2. Look up the timezone name based on the currency code.
	tzName, found := consts.CurrencyTimeZoneMap[currencyCode]
	if !found {
		// If the currency is not found in the map, return an error and the original UTC time.
		return currentUTCTime, fmt.Errorf("timezone information not found for currency '%s'", currencyCode)
	}

	// 3. Load the timezone location.
	//    time.LoadLocation loads timezone information from the system's or embedded timezone database.
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// If the timezone name is invalid or loading fails, return an error and the original UTC time.
		return currentUTCTime, fmt.Errorf("failed to load timezone '%s': %w", tzName, err)
	}

	// 4. Convert the current UTC time to the target timezone's time.
	convertedTime := currentUTCTime.In(loc)

	return convertedTime, nil
}

// ConvertTimeToCurrencyZone converts a given time.Time object to the timezone
// corresponding to the specified currency code.
//
// This function first ensures the input 'inputTime' is converted to UTC
// before performing the timezone conversion based on the currency.
//
// Parameters:
//
//	inputTime time.Time: The time.Time object to be converted. This can be
//	                     in any timezone (local, UTC, or specific).
//	currencyCode string: The currency code (e.g., "USD", "EUR").
//
// Returns:
//
//	time.Time: The converted time in the currency's timezone. If the currency
//	           is not found or the timezone fails to load, the time converted
//	           to UTC is returned.
//	error: An error if the currency's timezone information is not found
//	       or if the corresponding timezone fails to load.
func ConvertTimeToCurrencyZone(inputTime time.Time, currencyCode string) (time.Time, error) {
	// 1. First, ensure the input time is converted to UTC.
	// This makes the function robust regardless of the inputTime's original location.
	utcTime := inputTime.In(time.UTC)

	// 2. Look up the timezone name based on the currency code.
	tzName, found := consts.CurrencyTimeZoneMap[currencyCode]
	if !found {
		// If the currency isn't found, return an error and the UTC-normalized time.
		return utcTime, fmt.Errorf("timezone information not found for currency '%s'", currencyCode)
	}

	// 3. Load the timezone location.
	// time.LoadLocation fetches timezone data from the system's timezone database.
	// Ensure your deployment environment includes timezone data (e.g., tzdata package).
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// If the timezone name is invalid or loading fails, return an error and the UTC-normalized time.
		return utcTime, fmt.Errorf("failed to load timezone '%s' for currency '%s': %w", tzName, currencyCode, err)
	}

	// 4. Convert the UTC time to the target timezone's time.
	convertedTime := utcTime.In(loc)

	return convertedTime, nil
}

// ConvertTimeToCurrencyZoneStr converts a given time.Time object to the timezone
// corresponding to the specified currency code.
//
// This function first ensures the input 'inputTime' is converted to UTC
// before performing the timezone conversion based on the currency.
//
// Parameters:
//
//	inputTime time.Time: The time.Time object to be converted. This can be
//	                     in any timezone (local, UTC, or specific).
//	currencyCode string: The currency code (e.g., "USD", "EUR").
//
// Returns:
//
//	time.Time: The converted time in the currency's timezone. If the currency
//	           is not found or the timezone fails to load, the time converted
//	           to UTC is returned.
//	error: An error if the currency's timezone information is not found
//	       or if the corresponding timezone fails to load.
func ConvertTimeToCurrencyZoneStr(inputTime time.Time, currencyCode string) (string, error) {
	// 1. First, ensure the input time is converted to UTC.
	// This makes the function robust regardless of the inputTime's original location.
	utcTime := inputTime.In(time.UTC)

	// 2. Look up the timezone name based on the currency code.
	tzName, found := consts.CurrencyTimeZoneMap[currencyCode]
	if !found {
		// If the currency isn't found, return an error and the UTC-normalized time.
		return utcTime.Format(time.DateTime), fmt.Errorf("timezone information not found for currency '%s'", currencyCode)
	}

	// 3. Load the timezone location.
	// time.LoadLocation fetches timezone data from the system's timezone database.
	// Ensure your deployment environment includes timezone data (e.g., tzdata package).
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// If the timezone name is invalid or loading fails, return an error and the UTC-normalized time.
		return utcTime.Format(time.DateTime), fmt.Errorf("failed to load timezone '%s' for currency '%s': %w", tzName, currencyCode, err)
	}

	// 4. Convert the UTC time to the target timezone's time.
	convertedTime := utcTime.In(loc)

	return convertedTime.Format(time.DateTime), nil
}

// ConvertTimestampToCurrencyTime converts a given Unix timestamp to a time.Time object
// in the timezone corresponding to the specified currency.
//
// This function first converts the Unix timestamp to a UTC time.Time object.
// It then looks up and loads the appropriate timezone based on the currency code,
// finally converting the UTC time to that specific timezone.
//
// Parameters:
//
//	timestamp int64: The Unix timestamp (in seconds) to convert.
//	currencyCode string: The ISO 4217 currency code (e.g., "USD", "EUR").
//
// Returns:
//
//	time.Time: The converted time.Time object, with its Location() set to the
//	           currency's corresponding timezone. If the currency is not found
//	           or the timezone fails to load, the UTC time derived from the
//	           timestamp will be returned.
//	error: An error is returned if the currency's timezone information is not found
//	       or if the associated timezone fails to load.
func ConvertTimestampToCurrencyTime(timestamp int64, currencyCode string) (time.Time, error) {
	var seconds int64
	var nanoseconds int64

	// Convert the int64 timestamp to a string to check its length.
	timestampStr := strconv.FormatInt(timestamp, 10)
	length := len(timestampStr)
	// Determine precision based on string length heuristic.
	switch length {
	case 10: // e.g., 1748668009 (seconds)
		seconds = timestamp
		nanoseconds = 0
	case 13: // e.g., 1748668009000 (milliseconds)
		seconds = timestamp / 1000
		nanoseconds = (timestamp % 1000) * int64(time.Millisecond)
	case 16: // e.g., 1748668009000000 (microseconds)
		seconds = timestamp / 1_000_000
		nanoseconds = (timestamp % 1_000_000) * int64(time.Microsecond)
	case 19: // e.g., 1748668009000000000 (nanoseconds)
		seconds = timestamp / 1_000_000_000
		nanoseconds = timestamp % 1_000_000_000
	default:
		// For lengths not matching common patterns, it's ambiguous.
		// As a fallback, we'll try to treat it as milliseconds if it's large enough,
		// otherwise as seconds. A production system might return an error here
		// or require explicit precision.
		if timestamp > 1_000_000_000_000 { // Arbitrary threshold for large numbers
			seconds = timestamp / 1000
			nanoseconds = (timestamp % 1000) * int64(time.Millisecond)
		} else {
			seconds = timestamp
			nanoseconds = 0
		}
		// Return an error for unknown precision but still provide a best-effort time
		return time.Unix(seconds, nanoseconds).In(time.UTC), fmt.Errorf("ambiguous timestamp precision for '%d' (length %d)", timestamp, length)
	}
	// 1. Convert the Unix timestamp to a UTC time.Time object.
	//    time.Unix(sec, nsec) by default creates a time in the UTC timezone.
	// 1. Convert the Unix timestamp (now normalized to seconds and nanoseconds) to a UTC time.Time object.
	utcTime := time.Unix(seconds, nanoseconds).In(time.UTC) // Ensure it's explicitly UTC

	// 2. Look up the timezone name based on the currency code.
	tzName, found := consts.CurrencyTimeZoneMap[currencyCode] // Assuming consts.CurrencyTimeZoneMap is accessible
	if !found {
		// If the currency isn't found in the map, return an error and the UTC time derived from the timestamp.
		return utcTime, fmt.Errorf("timezone information not found for currency '%s'", currencyCode)
	}

	// 3. Load the timezone location.
	//    time.LoadLocation loads timezone information from the system's or embedded timezone database.
	//    Ensure your deployment environment (e.g., Docker container) includes timezone data
	//    (often provided by the 'tzdata' package) or this function might return an error.
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// If the timezone name is invalid or loading fails, return an error and the UTC time derived from the timestamp.
		return utcTime, fmt.Errorf("failed to load timezone '%s' for currency '%s': %w", tzName, currencyCode, err)
	}

	// 4. Convert the UTC time to the target timezone's time.
	convertedTime := utcTime.In(loc)

	return convertedTime, nil
}

// ConvertTimestampToCurrencyTimeStr converts a given Unix timestamp to a formatted
// date/time string in the timezone corresponding to the specified currency.
//
// This function first converts the Unix timestamp to a UTC time.Time object.
// It then looks up and loads the appropriate timezone based on the currency code,
// finally converting the UTC time to that specific timezone and formatting it.
//
// Parameters:
//
//	timestamp int64: The Unix timestamp (in seconds) to convert.
//	currencyCode string: The ISO 4217 currency code (e.g., "USD", "EUR").
//
// Returns:
//
//	string: The converted and formatted date/time string (YYYY-MM-DD HH:MM:SS).
//	        If the currency is not found or the timezone fails to load, the
//	        formatted UTC time derived from the timestamp will be returned.
//	error: An error is returned if the currency's timezone information is not found
//	       or if the associated timezone fails to load.
func ConvertTimestampToCurrencyTimeStr(timestamp int64, currencyCode string) (string, error) {

	var seconds int64
	var nanoseconds int64

	// Convert the int64 timestamp to a string to check its length.
	timestampStr := strconv.FormatInt(timestamp, 10)
	length := len(timestampStr)
	// Determine precision based on string length heuristic.
	switch length {
	case 10: // e.g., 1748668009 (seconds)
		seconds = timestamp
		nanoseconds = 0
	case 13: // e.g., 1748668009000 (milliseconds)
		seconds = timestamp / 1000
		nanoseconds = (timestamp % 1000) * int64(time.Millisecond)
	case 16: // e.g., 1748668009000000 (microseconds)
		seconds = timestamp / 1_000_000
		nanoseconds = (timestamp % 1_000_000) * int64(time.Microsecond)
	case 19: // e.g., 1748668009000000000 (nanoseconds)
		seconds = timestamp / 1_000_000_000
		nanoseconds = timestamp % 1_000_000_000
	default:
		// For lengths not matching common patterns, it's ambiguous.
		// As a fallback, we'll try to treat it as milliseconds if it's large enough,
		// otherwise as seconds. A production system might return an error here
		// or require explicit precision.
		if timestamp > 1_000_000_000_000 { // Arbitrary threshold for large numbers
			seconds = timestamp / 1000
			nanoseconds = (timestamp % 1000) * int64(time.Millisecond)
		} else {
			seconds = timestamp
			nanoseconds = 0
		}
		// Return an error for unknown precision but still provide a best-effort time
		return time.Unix(seconds, nanoseconds).In(time.UTC).Format(time.DateTime), fmt.Errorf("ambiguous timestamp precision for '%d' (length %d)", timestamp, length)
	}
	// 1. Convert the Unix timestamp to a UTC time.Time object.
	//    time.Unix(sec, nsec) by default creates a time in the UTC timezone.
	// 1. Convert the Unix timestamp (now normalized to seconds and nanoseconds) to a UTC time.Time object.
	utcTime := time.Unix(seconds, nanoseconds).In(time.UTC) // Ensure it's explicitly UTC

	// 2. Look up the timezone name based on the currency code.
	tzName, found := consts.CurrencyTimeZoneMap[currencyCode] // Assuming CurrencyTimeZoneMap is accessible
	if !found {
		// If the currency isn't found, return an error and the formatted UTC time.
		return utcTime.Format(time.DateTime), fmt.Errorf("timezone information not found for currency '%s'", currencyCode)
	}

	// 3. Load the timezone location.
	//    time.LoadLocation loads timezone information from the system's or embedded timezone database.
	//    Ensure your deployment environment (e.g., Docker container) includes timezone data
	//    (often provided by the 'tzdata' package) or this function might return an error.
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// If the timezone name is invalid or loading fails, return an error and the formatted UTC time.
		return utcTime.Format(time.DateTime), fmt.Errorf("failed to load timezone '%s' for currency '%s': %w", tzName, currencyCode, err)
	}

	// 4. Convert the UTC time to the target timezone's time.
	convertedTime := utcTime.In(loc)

	// 5. Format the converted time into the desired string layout.
	return convertedTime.Format(time.DateTime), nil
}

// ConvertUTCToCurrencyTime converts a given UTC time string
// to a time.Time object in the timezone corresponding to the specified currency.
//
// This function first parses the input 'utcTimeString' into a UTC time.Time object.
// It then looks up and loads the appropriate timezone based on the currency code,
// and finally converts the UTC time to that specific timezone.
//
// Parameters:
//
//	utcTimeString string: The UTC time string to convert (e.g., "2025-05-30T07:46:49Z" or "2025-05-30 07:46:49").
//	                      It's highly recommended to use RFC3339 format for robustness.
//	currencyCode string: The ISO 4217 currency code (e.g., "USD", "EUR").
//
// Returns:
//
//	time.Time: The converted time.Time object, with its Location() set to the
//	           currency's corresponding timezone. If parsing fails, currency is not found,
//	           or the timezone fails to load, a zero time.Time object will be returned.
//	error: An error is returned if string parsing fails, currency's timezone information
//	       is not found, or the associated timezone fails to load.
func ConvertUTCToCurrencyTime(utcTimeString string, currencyCode string) (time.Time, error) {
	var utcTime time.Time
	var err error

	// 1. Parse the input UTC time string into a time.Time object.
	// We'll try a few common UTC formats. For production, consider being more specific
	// about the expected format or use a robust parsing library if formats vary widely.
	layouts := []string{
		time.RFC3339,           // "2006-01-02T15:04:05Z" or "2006-01-02T15:04:05+00:00"
		"2006-01-02 15:04:05Z", // "YYYY-MM-DD HH:MM:SSZ"
		"2006-01-02 15:04:05",  // "YYYY-MM-DD HH:MM:SS" (assumed UTC if no zone info)
		// Add other possible UTC string formats you might receive
	}

	parsed := false
	for _, layout := range layouts {
		utcTime, err = time.Parse(layout, utcTimeString)
		if err == nil {
			// Ensure it's explicitly UTC after parsing.
			// If the layout didn't specify Z or +00:00, it might be parsed as local,
			// so force it to UTC.
			utcTime = utcTime.In(time.UTC)
			parsed = true
			break
		}
	}

	if !parsed {
		return time.Time{}, fmt.Errorf("failed to parse UTC time string '%s' with common layouts: %w", utcTimeString, err)
	}

	// 2. Look up the timezone name based on the currency code.
	tzName, found := consts.CurrencyTimeZoneMap[currencyCode] // Assuming CurrencyTimeZoneMap is accessible
	if !found {
		// If the currency isn't found in the map, return an error.
		return time.Time{}, fmt.Errorf("timezone information not found for currency '%s'", currencyCode)
	}

	// 3. Load the timezone location.
	//    time.LoadLocation loads timezone information from the system's or embedded timezone database.
	//    Ensure your deployment environment (e.g., Docker container) includes timezone data
	//    (often provided by the 'tzdata' package) or this function might return an error.
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// If the timezone name is invalid or loading fails, return an error.
		return time.Time{}, fmt.Errorf("failed to load timezone '%s' for currency '%s': %w", tzName, currencyCode, err)
	}

	// 4. Convert the UTC time to the target timezone's time.
	convertedTime := utcTime.In(loc)

	return convertedTime, nil
}

// GetMondayFirstWeekRange 获取给定日期所在周的开始和结束日期，以周一为第一天。
func GetMondayFirstWeekRange(t time.Time) (start, end time.Time) {
	// 获取当前日期是星期几
	// Weekday() 返回值：Sunday=0, Monday=1, ..., Saturday=6
	weekday := t.Weekday()

	// 计算需要倒退的天数，以周一为起点
	// (weekday - time.Monday + 7) % 7 确保结果为 0-6 之间的正数
	// 例如：今天是周三 (3)， (3 - 1 + 7) % 7 = 2，需要倒退2天
	// 例如：今天是周一 (1)， (1 - 1 + 7) % 7 = 0，需要倒退0天
	offset := (weekday - time.Monday + 7) % 7

	// 找到本周开始日期（周一），并将时间重置为当天的 00:00:00
	start = t.AddDate(0, 0, -int(offset))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	// 找到本周结束日期（周日），并将时间重置为当天的 23:59:59
	end = start.AddDate(0, 0, 6)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

	return start, end
}

// GetSundayFirstWeekRange 获取给定日期所在周的开始和结束日期，以周日为第一天。
func GetSundayFirstWeekRange(t time.Time) (start, end time.Time) {
	// 获取当前日期是星期几
	// Weekday() 返回值：Sunday=0, Monday=1, ..., Saturday=6
	weekday := t.Weekday()

	// 找到本周开始日期（周日），并将时间重置为当天的 00:00:00
	start = t.AddDate(0, 0, -int(weekday))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	// 找到本周结束日期（周六），并将时间重置为当天的 23:59:59
	end = start.AddDate(0, 0, 6)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

	return start, end
}

// IsTimeBetween 判断一个时间点 t 是否在 start 和 end 之间（包含边界）。
//
// 注意：此函数在比较时会自动处理时间戳和时区差异，所以无需手动转换为 UTC 时间。
// 参数:
//
//	t: 需要判断的时间
//	start: 时间段的开始时间
//	end: 时间段的结束时间
//
// 返回值:
//
//	如果 t 在 [start, end] 范围内，返回 true；否则返回 false。
func IsTimeBetween(t, start, end time.Time) bool {
	// 关键：t.After()、t.Before() 和 t.Equal() 方法在内部已经处理了时区问题，
	// 它们会基于时间点（而不是时区）进行精确比较。
	return (t.After(start) || t.Equal(start)) && (t.Before(end) || t.Equal(end))
}

func GenerateTimeRanges(start, end time.Time, intervalType int32) ([]time.Time, error) {

	if start.After(end) {
		return nil, fmt.Errorf("start time no later than end time")
	}

	var interval time.Duration
	switch intervalType {
	case 1:
		interval = time.Minute
	case 2:
		interval = 5 * time.Minute
	case 3:
		interval = 30 * time.Minute
	case 4:
		interval = time.Hour
	case 5:
		interval = time.Hour * 24
	default:
		return nil, fmt.Errorf("Unsupported time interval type : %d", intervalType)
	}

	var timeRanges []time.Time
	current := start

	for !current.After(end) {
		timeRanges = append(timeRanges, current)
		current = current.Add(interval)
	}

	return timeRanges, nil
}
