package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseTime parses a time string in the format "hours:minutes:seconds"
// where each component (hours, minutes, seconds) is optional.
// If hours is not provided, it defaults to 0.
// If only minutes are provided, it is interpreted as minutes and seconds.
// If only seconds are provided, it is interpreted as seconds.
// The function returns a time.Duration representing the parsed time.
// An error is returned if the input string is in an invalid format.

func ParseTime(input string) (time.Duration, error) {
	if len(input) == 0 {
		return 0, fmt.Errorf("invalid time format")
	}
	parts := strings.Split(input, ":")
	var hours, minutes, seconds int

	switch len(parts) {
	case 1:
		seconds, _ = strconv.Atoi(parts[0])
	case 2:
		minutes, _ = strconv.Atoi(parts[0])
		seconds, _ = strconv.Atoi(parts[1])
	case 3:
		hours, _ = strconv.Atoi(parts[0])
		minutes, _ = strconv.Atoi(parts[1])
		seconds, _ = strconv.Atoi(parts[2])
	default:
		return 0, fmt.Errorf("invalid time format")
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return duration, nil
}
